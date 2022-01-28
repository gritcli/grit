package commands

import (
	"context"
	"errors"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/render"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

// newCloneCommand returns the "clone" command.
func newCloneCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <repo>",
		Args:  cobra.ExactArgs(1),
		Short: "Clone a remote repository",
		Long: heredoc.Doc(`
		The "clone" command makes a local clone of a remote repository then
		changes the users's current working directory to that of the clone.

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
				remote, err := resolveRemoteRepo(
					ctx,
					cmd,
					client,
					clientOptions,
					args[0],
				)
				if err != nil {
					return err
				}

				local, err := cloneRepo(
					ctx,
					cmd,
					client,
					clientOptions,
					remote,
				)
				if err != nil {
					return err
				}

				cmd.Println(render.RelPath(local.AbsoluteCloneDir))

				return executor("cd", local.AbsoluteCloneDir)
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
				Word:   word,
				Filter: api.SuggestReposFilter_SUGGEST_REMOTE_ONLY,
			})
		}),
	}
}
