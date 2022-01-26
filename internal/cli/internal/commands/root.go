package commands

import (
	"github.com/gritcli/grit/internal/cli/internal/deps"
	"github.com/gritcli/grit/internal/cli/internal/flags"
	"github.com/spf13/cobra"
)

// NewRoot returns the root command.
//
// v is the version to display. It is passed from the main package where it is
// made available as part of the build process.
func NewRoot(v string) *cobra.Command {
	root := &cobra.Command{
		Version:      v,
		Use:          "grit",
		Short:        "Manage your local source repository clones",
		SilenceUsage: true, // otherwise ANY error shows usage
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Add the currently-executing Cobra CLI command the the DI
			// container.
			//
			// This hook is called after the CLI arguments are resolved to a
			// specific command.
			//
			// This allows other DI provider definitions to make use of the
			// flags passed to the command.
			deps.Container.Provide(func() *cobra.Command {
				return cmd
			})
		},
	}

	flags.SetupVerbose(root)
	flags.SetupNoInteractive(root)
	flags.SetupSocket(root)
	flags.SetupShellExecutorOutput(root)

	root.AddCommand(
		newCloneCommand(),
		newShellIntegrationCommand(),
		newSourceCommand(),
	)

	return root
}
