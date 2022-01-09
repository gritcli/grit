package deps

import (
	"github.com/gritcli/grit/internal/cli/internal/flags"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/gritcli/grit/internal/common/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
	) (grpc.ClientConnInterface, error) {
		conn, err := grpc.Dial(
			"unix:"+cfg.Daemon.Socket,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)
		if err != nil {
			return nil, err
		}
		Container.Defer(conn.Close)

		return conn, nil
	})

	Container.Provide(func(
		cmd *cobra.Command,
	) *api.ClientOptions {
		return &api.ClientOptions{
			CaptureDebugLog: flags.IsVerbose(cmd),
		}
	})

	Container.Provide(api.NewAPIClient)
}
