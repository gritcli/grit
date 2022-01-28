package interactive

import (
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// SelectRemoteRepos prompts the user to select from a list of repositories.
func SelectRemoteRepos(
	cmd *cobra.Command,
	repos []*api.RemoteRepo,
) (*api.RemoteRepo, error) {
	n, err := selectPrompt(
		cmd,
		"Chhose a repository:",
		repos,
		&promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   `> {{ .Name | cyan }} {{ print "[" .Source "]" | yellow }} {{ .Description         }}  {{ .WebUrl }}`,
			Inactive: `  {{ .Name | faint }} {{ print "[" .Source "]" | faint }} {{ .Description | faint }}  {{ .WebUrl | faint }}`,
		},
		func(s string, i int) bool {
			s = strings.ToLower(s)

			if strings.Contains(strings.ToLower(repos[i].Name), s) {
				return true
			}

			if strings.Contains(strings.ToLower(repos[i].Source), s) {
				return true
			}

			if strings.Contains(strings.ToLower(repos[i].Description), s) {
				return true
			}

			return false
		},
	)

	if err != nil {
		return nil, err
	}

	return repos[n], nil
}

// selectPrompt prompts the user to choose from a set of options.
func selectPrompt(
	cmd *cobra.Command,
	label string,
	items interface{},
	templates *promptui.SelectTemplates,
	searchPredicate func(string, int) bool,
) (int, error) {
	n := reflect.ValueOf(items).Len()

	if n == 0 {
		panic("no options provided")
	}

	if n == 1 {
		return 0, nil
	}

	if !flags.IsInteractive(cmd) {
		return 0, errors.New("can not prompt to select an option, command is non-interactive")
	}

	p := promptui.Select{
		Label:        label,
		Items:        items,
		Templates:    templates,
		HideSelected: true,
		Searcher:     searchPredicate,
		Stdin:        cmd.InOrStdin().(io.ReadCloser),
		Stdout:       cmd.OutOrStdout().(io.WriteCloser),
	}

	n, _, err := p.Run()
	return n, err
}
