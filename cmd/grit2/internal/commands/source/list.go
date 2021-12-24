package source

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/internal/di/cobradi"
	"github.com/gritcli/grit/internal/source"
	"github.com/spf13/cobra"
)

// newListCommand returns the "source ls" command.
func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Short:   "List the configured repository sources",
		Long: heredoc.Doc(`
		The "source list" command lists the configured repository sources and
		their current status.
		`),
		RunE: cobradi.RunE(func(
			cmd *cobra.Command,
			args []string,
			sources []source.Source,
		) error {
			ctx := cmd.Context()

			for _, src := range sources {
				status, err := src.Status(ctx)
				if err != nil {
					status = "error: " + err.Error()
				}

				cmd.Printf(
					"%s: %s\n",
					src.Description(),
					status,
				)
			}

			return nil
		}),
	}
}
