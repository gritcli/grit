package interactive

import (
	"errors"
	"io"
	"reflect"

	"github.com/gritcli/grit/internal/client/internal/flags"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// selectPrompt prompts the user to choose from a set of options.
func selectPrompt(
	cmd *cobra.Command,
	label string,
	items interface{},
	templates *promptui.SelectTemplates,
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
		Label:     label,
		Items:     items,
		Templates: templates,
		Stdin:     cmd.InOrStdin().(io.ReadCloser),
		Stdout:    cmd.OutOrStdout().(io.WriteCloser),
	}

	n, _, err := p.Run()
	return n, err
}

// SelectRepos prompts the user to select from a list of repositories.
func SelectRepos(
	cmd *cobra.Command,
	label string,
	repos []*api.Repo,
) (*api.Repo, error) {
	n, err := selectPrompt(
		cmd,
		label,
		repos,
		&promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   `> {{ .Source         }}: {{ .Name | cyan  }}  {{ .Description         }}  {{ .WebUrl }}`,
			Inactive: `  {{ .Source | faint }}: {{ .Name | faint }}  {{ .Description | faint }}  {{ .WebUrl | faint }}`,
			Selected: `  {{ .Source | faint }}: {{ .Name | green }}  {{ .Description | faint }}  {{ .WebUrl | faint }}`,
		},
	)

	if err != nil {
		return nil, err
	}

	return repos[n], nil
}
