package clone

import (
	"context"
	_ "embed"
	"errors"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/completion"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/gritcli/grit/cli/internal/render"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

//go:embed help.txt
var helpText string

// Command returns the "clone" command.
func Command(con *imbue.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "clone [--from-source <source> [--no-resolve]] <repo>",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Short:                 "Clone a remote repository",
		Long:                  helpText,
		ValidArgsFunction: completion.Positional(
			completion.RepoName(con, api.Locality_REMOTE),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOrID := args[0]

			if queryOrID == "" {
				return errors.New("<repo> argument must not be empty")
			}

			source, noResolve, err := flags.FromSource(cmd)
			if err != nil {
				return err
			}

			cmd.SilenceUsage = true

			return imbue.Invoke3(
				cmd.Context(),
				con,
				func(
					ctx context.Context,
					client api.APIClient,
					options *api.ClientOptions,
					exec shell.Executor,
				) error {
					if !noResolve {
						repo, ok, err := resolve(
							ctx,
							cmd,
							client,
							options,
							queryOrID,
							source,
						)
						if err != nil {
							return err
						}
						if !ok {
							return nil
						}

						queryOrID = repo.GetId()
						source = repo.GetSource()
					}

					dir, err := clone(
						ctx,
						cmd,
						client,
						options,
						queryOrID,
						source,
					)
					if err != nil {
						return err
					}

					cmd.Println(render.RelPath(dir))

					return exec("cd", dir)
				},
			)
		},
	}

	flags.SetupFromSource(cmd, con)

	return cmd
}
