package source

import (
	"github.com/jmalloc/grit/cmd/grit2/internal/commands/source/add"
	"github.com/spf13/cobra"
)

// NewRoot returns the "source" command.
func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "manage repository sources",
	}

	cmd.AddCommand(
		add.NewRoot(),
		newListCommand(),
	)

	return cmd
}
