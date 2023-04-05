package source

import (
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/commands/source/list"
	"github.com/gritcli/grit/cli/internal/commands/source/login"
	"github.com/spf13/cobra"
)

// Command returns the "source" command.
func Command(con *imbue.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "source",
		DisableFlagsInUseLine: true,
		Short:                 "Manage repository sources",
	}

	cmd.AddCommand(
		list.Command(con),
		login.Command(con),
	)

	return cmd
}
