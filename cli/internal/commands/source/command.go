package source

import (
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/commands/source/list"
	"github.com/spf13/cobra"
)

// Command returns the "source" command.
func Command(c *imbue.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Manage repository sources",
	}

	cmd.AddCommand(
		list.Command(c),
	)

	return cmd
}
