package commands

import (
	"os"
	"path/filepath"

	"github.com/gritcli/grit/cmd/grit2/internal/commands/source"
	"github.com/gritcli/grit/internal/di"
	"github.com/gritcli/grit/internal/di/cobradi"
	"github.com/spf13/cobra"
)

// NewRoot returns the root command.
//
// v is the version to display. It is passed from the main package where it is
// made available as part of the build process.
func NewRoot(v string) *cobra.Command {
	root := &cobra.Command{
		Version: v,
		Use:     executableName(),
		Short:   "keep track of your local git clones",
		PersistentPreRunE: cobradi.Setup(
			func(
				c *di.Container,
				cmd *cobra.Command,
			) {
				provideConfig(c, cmd)
				provideShellExecutor(c, cmd)
			},
		),
	}

	setupConfig(root)
	setupShellExecutor(root)

	root.AddCommand(
		source.NewRoot(),
		newShellIntegrationCommand(),
		newCloneCommand(),
	)

	return root
}

// executableName returns the name of the grit executable.
func executableName() string {
	return filepath.Base(os.Args[0])
}
