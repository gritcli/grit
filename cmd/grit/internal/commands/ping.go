package commands

import (
	"context"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/cmd/grit/internal/deps"
	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/commondeps"
	"github.com/spf13/cobra"
)

// newPingCommand returns the "ping" command.
func newPingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ping",
		Args:  cobra.NoArgs,
		Short: "Ping the Grit daemon",
		Long: heredoc.Doc(`
		The "ping" checks that the Grit CLI can communicate with the Grit daemon.
		`),
		RunE: deps.Run(func(
			ctx context.Context,
			cmd *cobra.Command,
			ping api.PingClient,
			info commondeps.ExecutableInfo,
		) error {
			res, err := ping.Ping(ctx, &api.PingRequest{
				Version: info.Version,
			})
			if err != nil {
				return err
			}

			cmd.Printf("Daemon version: %s\n", res.Version)
			cmd.Printf("CLI version: %s", info.Version)

			return nil
		}),
	}
}
