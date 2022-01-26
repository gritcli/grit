package commands

import (
	"context"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/internal/cli/internal/deps"
	"github.com/gritcli/grit/internal/cli/internal/render"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/spf13/cobra"
)

// newSourceListCommand returns the "source list" command.
func newSourceListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Short:   "List the configured repository sources",
		Long: heredoc.Doc(`
		The "source list" command lists the configured repository sources and
		their current status.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deps.Invoke(cmd, func(
				ctx context.Context,
				client api.APIClient,
			) error {
				res, err := client.Sources(ctx, &api.SourcesRequest{})
				if err != nil {
					return err
				}

				for _, src := range res.Sources {
					cmd.Printf(
						"%s\t%s\t%s\t%s\n",
						src.Name,
						src.Description,
						src.Status,
						render.AbsPath(src.CloneDir),
					)
				}

				return nil
			})
		},
	}
}
