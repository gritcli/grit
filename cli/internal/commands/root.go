package commands

import (
	"os"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/commands/clone"
	"github.com/gritcli/grit/cli/internal/commands/setupshell"
	"github.com/gritcli/grit/cli/internal/commands/source"
	"github.com/gritcli/grit/cli/internal/commands/version"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/spf13/cobra"
)

// Root returns the root "grit" command.
func Root(con *imbue.Container, ver string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "grit",
		DisableFlagsInUseLine: true,
		Short:                 "Manage your local VCS clones",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Add the currently-executing Cobra CLI command the the DI
			// container.
			//
			// This hook is called after the CLI arguments are resolved to a
			// specific command.
			//
			// This allows other DI provider definitions to make use of the
			// flags passed to the command.
			imbue.With0(
				con,
				func(
					ctx imbue.Context,
				) (*cobra.Command, error) {
					return cmd, nil
				},
			)
		},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	flags.SetupVerbose(cmd)
	flags.SetupNoInteractive(cmd)
	flags.SetupSocket(cmd)
	flags.SetupShellExecutorOutput(cmd)

	cmd.AddCommand(
		clone.Command(con),
		setupshell.Command(con),
		source.Command(con),
		version.Command(con, ver),
	)

	return cmd
}
