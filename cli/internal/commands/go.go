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

// newGoCommand returns the "go" command.
func newGoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "go <repo>",
		Args:  cobra.ExactArgs(1),
		Short: "Change the working directory to a repository, cloning if necessary",
		Long: heredoc.Doc(`
		The "go" command changes the current working directory to that of an
		existing local clone, or clones a remote repository if none is found.

		It is a streamlined alternative to running the "chdir" command, followed
		by "clone" if changing directories fails.

		The <repo> argument is a repository name (or a part thereof), URL, or
		other identifier. For example, the Grit repository itself may be
		referred to as "gritcli/grit" or simply "grit".

		If there are multiple matching local clones and the shell is interactive
		the user is prompted to select the desired repository.
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

			return client.SuggestRepo(ctx, &api.SuggestRepoRequest{
				Word:          word,
				IncludeLocal:  true,
				IncludeRemote: true,
			})
		}),
	}
}
