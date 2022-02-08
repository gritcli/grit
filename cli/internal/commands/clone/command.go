package clone

import (
	"context"
	_ "embed"
	"errors"

	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/completion"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/gritcli/grit/cli/internal/render"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

//go:embed help.txt
var helpText string

// Command is the "clone" command.
var Command = &cobra.Command{
	Use:               "clone [--from-source <source> [--no-resolve]] <repo>",
	Args:              cobra.ExactArgs(1),
	Short:             "Clone a remote repository",
	Long:              helpText,
	ValidArgsFunction: completion.RepoName(0, api.Locality_REMOTE),
	RunE:              run,
}

func init() {
	flags.SetupFromSource(Command)
}

// run executes the command.
func run(cmd *cobra.Command, args []string) error {
	queryOrID := args[0]

	if queryOrID == "" {
		return errors.New("<repo> argument must not be empty")
	}

	source, noResolve, err := flags.FromSource(cmd)
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	return cobradi.Invoke(cmd, func(
		ctx context.Context,
		client api.APIClient,
		options *api.ClientOptions,
		exec shell.Executor,
	) error {
		if !noResolve {
			var err error
			queryOrID, source, err = resolve(
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
	})
}
