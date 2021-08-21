package source

import (
	"github.com/gritcli/grit/cmd/grit2/internal/commands/source/setup"
	"github.com/spf13/cobra"
)

// NewRoot returns the "source" command.
func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "manage repository sources",
	}

	cmd.AddCommand(
		setup.NewRoot(),
		newListCommand(),
	)

	return cmd
}
