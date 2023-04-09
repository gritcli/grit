package source

import (
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/commands/source/list"
	"github.com/gritcli/grit/cli/internal/commands/source/signin"
	"github.com/gritcli/grit/cli/internal/commands/source/signout"
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
		signin.Command(con),
		signout.Command(con),
	)

	return cmd
}
