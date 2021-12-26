package deps

import (
	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	Container.Provide(func(cfg config.Config) (grpc.ClientConnInterface, error) {
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

	Container.Provide(api.NewAPIClient)
}
