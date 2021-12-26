package source

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/cmd/grit/internal/deps"
	"github.com/gritcli/grit/internal/api"
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
		RunE: deps.Run(func(
			cmd *cobra.Command,
			args []string,
			client api.APIClient,
		) error {
			ctx := cmd.Context()

			res, err := client.ListSources(ctx, &api.ListSourcesRequest{})
			if err != nil {
				return err
			}

			for _, src := range res.Sources {
				cmd.Printf(
					"%s (%s): %s\n",
					src.Name,
					src.Driver,
					src.Config,
				)
			}

			return nil
		}),
	}
}
