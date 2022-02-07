package source

import (
	"github.com/gritcli/grit/cli/internal/commands/source/list"
	"github.com/spf13/cobra"
)

// Command is the "source" command.
var Command = &cobra.Command{
	Use:   "source",
	Short: "Manage repository sources",
}

func init() {
	Command.AddCommand(
		list.Command,
	)
}
