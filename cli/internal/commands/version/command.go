package version

import (
	"context"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/spf13/cobra"
)

// Command returns the "version" command.
func Command(con *imbue.Container, ver string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Short:                 "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return imbue.Invoke1(
				cmd.Context(),
				con,
				func(
					ctx context.Context,
					client api.APIClient,
				) error {
					cmd.SilenceUsage = true

					cmd.Printf(
						"grit cli version\t%s\n",
						ver,
					)

					res, err := client.DaemonInfo(ctx, &api.DaemonInfoRequest{})
					if err != nil {
						return err
					}

					cmd.Printf(
						"grit daemon version\t%s\n",
						res.GetVersion(),
					)

					return nil
				},
			)
		},
	}

	return cmd
}
