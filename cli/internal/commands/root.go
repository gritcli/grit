package commands

import (
	"os"

	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/commands/clone"
	"github.com/gritcli/grit/cli/internal/commands/source"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/spf13/cobra"
)

// Root is the root "grit" command.
var Root = &cobra.Command{
	Use:   "grit2",
	Short: "Manage your local VCS clones",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Add the currently-executing Cobra CLI command the the DI
		// container.
		//
		// This hook is called after the CLI arguments are resolved to a
		// specific command.
		//
		// This allows other DI provider definitions to make use of the
		// flags passed to the command.
		cobradi.Provide(cmd, func() *cobra.Command {
			return cmd
		})
	},
}

func init() {
	// Explicitly set the streams for Cobra to use; otherwise, cmd.Print() uses
	// stderr.
	Root.SetIn(os.Stdin)
	Root.SetOut(os.Stdout)
	Root.SetErr(os.Stderr)

	flags.SetupVerbose(Root)
	flags.SetupNoInteractive(Root)
	flags.SetupSocket(Root)
	flags.SetupShellExecutorOutput(Root)

	Root.AddCommand(
		clone.Command,
		source.Command,
	)
}
