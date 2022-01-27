package deps

import (
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	Container.Provide(func(
		cmd *cobra.Command,
	) (grpc.ClientConnInterface, error) {
		socket, err := flags.Socket(cmd)
		if err != nil {
			return nil, err
		}

		conn, err := grpc.Dial(
			"unix:"+socket,
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
