package cli

import (
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/api/daemonapi"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	imbue.With1(
		catalog,
		func(
			ctx imbue.Context,
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
			ctx.Defer(conn.Close)

			return conn, nil
		},
	)

	imbue.With1(
		catalog,
		func(
			ctx imbue.Context,
			cmd *cobra.Command,
		) (*api.ClientOptions, error) {
			return &api.ClientOptions{
				Verbose: flags.IsVerbose(cmd),
			}, nil
		},
	)

	imbue.With1(
		catalog,
		func(
			ctx imbue.Context,
			conn grpc.ClientConnInterface,
		) (api.APIClient, error) {
			return api.NewAPIClient(conn), nil
		},
	)

	imbue.With1(
		catalog,
		func(
			ctx imbue.Context,
			conn grpc.ClientConnInterface,
		) (daemonapi.APIClient, error) {
			return daemonapi.NewAPIClient(conn), nil
		},
	)
}
