package deps

import (
	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/di"
	"google.golang.org/grpc"
)

func init() {
	Container.Provide(func(x di.ExecutableInfo) api.PingServer {
		return &apiserver.PingServer{
			Version: x.Version,
		}
	})

	Container.Provide(func(
		ping api.PingServer,
	) *grpc.Server {
		s := grpc.NewServer()

		api.RegisterPingServer(s, ping)

		return s
	})

}
