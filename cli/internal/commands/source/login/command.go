package login

import (
	"context"
	_ "embed"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/completion"
	"github.com/gritcli/grit/cli/internal/render"
	"github.com/spf13/cobra"
)

//go:embed help.txt
var helpText string

// Command returns the "source auth" command.
func Command(con *imbue.Container) *cobra.Command {
	return &cobra.Command{
		Use:                   "login",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Aliases:               []string{"auth"},
		Short:                 "Authenticate with a specific repository source",
		Long:                  helpText,
		ValidArgsFunction: completion.Positional(
			completion.SourceName(con),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return imbue.Invoke1(
				cmd.Context(),
				con,
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
