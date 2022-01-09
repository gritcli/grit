package commands

import (
	"github.com/spf13/cobra"
)

// newSourceCommand returns the "source" command.
func newSourceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Manage repository sources",
	}

	cmd.AddCommand(
		newSourceListCommand(),
	)

	return cmd
}
