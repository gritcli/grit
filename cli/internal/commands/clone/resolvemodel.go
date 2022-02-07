package clone

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/apitea"
)

// resolutionComplete is a tea.Msg that indicates resolution has completed.
type resolutionComplete struct {
	// Repo is the selected repo. It is nil if the user aborted the process.
	Repo *api.RemoteRepo

	// Error is the error that prevented the resolution process from succeeding.
	Error error
}

// resolveModel is the model used for interactive repository resolution.
type resolveModel struct {
	Repo  *api.RemoteRepo
	Error error

	query     string
	responses api.API_ResolveRepoClient
	spinner   spinner.Model
	output    []string
	matches   []*api.RemoteRepo
	cursor    int
	done      bool
}

// newResolveModel returns a new model for interactive repository resolution.
func newResolveModel(
	query string,
	responses api.API_ResolveRepoClient,
) resolveModel {
	m := resolveModel{
		query:     query,
		responses: responses,
		spinner:   spinner.New(),
	}

	m.spinner.Spinner = spinner.Points

	return m
}

func (m resolveModel) Init() tea.Cmd {
	return tea.Batch(
		apitea.WaitForResolveResponse(m.responses, ""),
		m.spinner.Tick,
	)
}

func (m resolveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var batch []tea.Cmd

	switch msg := msg.(type) {
	case resolutionComplete:
		m.Repo = msg.Repo
		m.Error = msg.Error
		return m, tea.Quit

	case tea.KeyMsg:
		if cmd := m.handleKeyPress(msg); cmd != nil {
			batch = append(batch, cmd)
		}

	case apitea.OutputReceived:
		m.output = append(m.output, msg.Output.GetMessage())
		batch = append(batch, apitea.WaitForResolveResponse(m.responses, ""))

	case apitea.RemoteRepoMatched:
		m.matches = append(m.matches, msg.Repo)
		batch = append(batch, apitea.WaitForResolveResponse(m.responses, ""))

	case apitea.ResolveComplete:
		if cmd := m.handleComplete(); cmd != nil {
			batch = append(batch, cmd)
		}

	case apitea.ResolveFailed:
		return m, func() tea.Msg {
			return resolutionComplete{
				Error: msg.Error,
			}
		}
	}

	if !m.done {
		s, cmd := m.spinner.Update(msg)
		m.spinner = s
		batch = append(batch, cmd)
	}

	return m, tea.Batch(batch...)
}

// handleKeyProcess updates the model as a result of a key being pressed.
func (m *resolveModel) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "esc":
		return func() tea.Msg {
			return resolutionComplete{}
		}

	case "up", "j":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "k":
		if m.cursor < len(m.matches)-1 {
			m.cursor++
		}

	case "enter":
		if len(m.matches) > 0 {
			return func() tea.Msg {
				return resolutionComplete{
					Repo: m.matches[m.cursor],
				}
			}
		}
	}

	return nil
}

// handleComplete handles completion of the API resolution operation.
func (m *resolveModel) handleComplete() tea.Cmd {
	m.done = true

	switch len(m.matches) {
	case 0:
		return func() tea.Msg {
			return resolutionComplete{
				Error: errors.New("no matching repositories"),
			}
		}

	case 1:
		return func() tea.Msg {
			return resolutionComplete{
				Repo: m.matches[0],
			}
		}
	}

	return nil
}

func (m resolveModel) View() string {
	count := len(m.matches)
	view := ""

	if len(m.output) > 0 {
		for _, o := range m.output {
			view += o + "\n"
		}

		view += "\n"
	}

	view += fmt.Sprintf("resolving '%s', ", m.query)

	switch count {
	case 0:
		view += "no matching repositories found"
	case 1:
		view += "one matching repository found"
	default:
		view += fmt.Sprintf("%d matching repositories found", len(m.matches))
	}

	if !m.done {
		view += " " + m.spinner.View()
	}

	view += "\n"

	selected := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))
	faint := lipgloss.NewStyle().Faint(true)

	switch count {
	case 0:
		if m.done {
			view += faint.Render("nothing to clone")
		} else {
			view += faint.Render("press [esc] to cancel")
		}

	case 1:
		if m.done {
			view += faint.Render(
				fmt.Sprintf("chose '%s' automatically",
					m.matches[0].GetName(),
				),
			)
		} else {
			view += faint.Render("press [enter] to clone ") +
				selected.Render(m.matches[0].GetName()) +
				faint.Render(" immediately, or [esc] to cancel")
		}
	default:
		view += faint.Render(
			"use [↑/↓] to select a repository, [enter] to clone, or [esc] to cancel",
		)
	}

	if !m.done {
		view += faint.Render(" (or wait for more matches)")
	}

	view += "\n\n"

	if count > 1 {
		for i, r := range m.matches {
			if i == m.cursor {
				view += " ∙ "
				view += selected.Render(r.GetName())
			} else {
				view += "   "
				view += r.GetName()
			}

			view += " (" + r.GetSource() + ")" + "\n"
			view += "   "
			view += faint.Render(r.GetDescription()) + "\n"
			view += "\n"
		}
	}

	return view
}
