package list

import (
	"context"
	_ "embed"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/render"
	"github.com/spf13/cobra"
)

//go:embed help.txt
var helpText string

// Command returns the "source list" command.
func Command(container *imbue.Container) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Short:   "List the configured repository sources",
		Long:    helpText,
		RunE: func(cmd *cobra.Command, args []string) error {
			return imbue.Invoke1(
				cmd.Context(),
				container,
				func(
					ctx context.Context,
					client api.APIClient,
				) error {
					res, err := client.ListSources(ctx, &api.ListSourcesRequest{})
					if err != nil {
						return err
					}

					for _, src := range res.Sources {
						cmd.Printf(
							"%s\t%s\t%s\t%s\n",
							src.Name,
							src.Description,
							src.Status,
							render.AbsPath(src.BaseCloneDir),
						)
					}

					return nil
				},
			)
		},
	}
}
