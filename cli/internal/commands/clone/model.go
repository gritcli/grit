package clone

import (
	"errors"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gritcli/grit/api"
)

type model struct {
	Stream     api.API_ResolveRepoClient
	ChosenRepo *api.RemoteRepo
	Error      error

	repos    []*api.RemoteRepo
	selected int
	output   string
}

func (m model) Init() tea.Cmd {
	return m.waitForRepo
}

func (m model) waitForRepo() tea.Msg {
	res, err := m.Stream.Recv()
	if err != nil {
		return err
	}

	return res
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		case "up", "j":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "k":
			if m.selected < len(m.repos)-1 {
				m.selected++
			}
		case "enter":
			if len(m.repos) != 0 {
				m.ChosenRepo = m.repos[m.selected]
				return m, tea.Quit
			}
		}

	case *api.ResolveRepoResponse:
		if out := msg.GetOutput(); out != nil {
			m.output += out.Message + "\n"
		} else if r := msg.GetRemoteRepo(); r != nil {
			m.repos = append(m.repos, r)
		}

		return m, m.waitForRepo

	case error:
		if msg != io.EOF {
			m.Error = msg
			return m, tea.Quit
		}

		if len(m.repos) == 1 {
			m.ChosenRepo = m.repos[0]
			return m, tea.Quit
		}

		if len(m.repos) == 0 {
			m.Error = errors.New("no matching repositories found")
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	view := m.output

	if view != "" {
		view += "\n"
	}

	for i, r := range m.repos {
		item := "  "
		if i == m.selected {
			item = "> "
		}

		item += fmt.Sprintf(
			"%s %s",
			r.GetName(),
			r.GetDescription(),
		)

		view += item + "\n"
	}

	return view
}
