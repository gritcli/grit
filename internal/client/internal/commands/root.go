package commands

import (
	"context"

	"github.com/gritcli/grit/internal/client/internal/commands/source"
	"github.com/gritcli/grit/internal/client/internal/deps"
	"github.com/gritcli/grit/internal/client/internal/flags"
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
			// container. This hook is called after the CLI arguments are
			// resolved to a specific command.
			deps.Container.Provide(func() (
				context.Context,
				*cobra.Command,
				[]string,
			) {
				return cmd.Context(), cmd, args
			})
		},
	}

	flags.SetupVerbose(root)
	flags.SetupNoInteractive(root)
	flags.SetupConfig(root)
	flags.SetupShellExecutorOutput(root)

	root.AddCommand(
		source.NewRoot(),
		newShellIntegrationCommand(),
		newCloneCommand(),
	)

	return root
}
