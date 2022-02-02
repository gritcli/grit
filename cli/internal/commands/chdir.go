package commands

import (
	"context"
	"errors"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

// newChDirCommand returns the "chdir" command.
func newChDirCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "chdir <repo>",
		Aliases: []string{"cd"},
		Args:    cobra.ExactArgs(1),
		Short:   "Change the working directory to an existing local clone",
		Long: heredoc.Doc(`
		The "chdir" command changes the current working directory to that of an
		existing local clone.

		The <repo> argument is a repository name (or a part thereof), URL, or
		other identifier. For example, the Grit repository itself may be
		referred to as "gritcli/grit" or simply "grit".

		If there are multiple matches and the shell is interactive the user is
		prompted to select the desired repository.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("<repo> argument must not be empty")
			}

			return cobradi.Invoke(cmd, func(
				ctx context.Context,
				client api.APIClient,
				clientOptions *api.ClientOptions,
				executor shell.Executor,
			) error {
				return nil
			})
		},
		ValidArgsFunction: suggest(func(
			ctx context.Context,
			client api.APIClient,
			cmd *cobra.Command,
			args []string,
			word string,
		) (*api.SuggestResponse, error) {
			if len(args) != 0 {
				return nil, nil
			}

			return client.SuggestRepos(ctx, &api.SuggestReposRequest{
				Word:     word,
				Locality: api.Locality_LOCAL_ONLY,
			})
		}),
	}
}
