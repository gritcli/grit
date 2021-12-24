package deps

import (
	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/commondeps"
	"google.golang.org/grpc"
)

func init() {
	Container.Provide(func(
		info commondeps.ExecutableInfo,
	) api.PingServer {
		return &apiserver.PingServer{
			Version: info.Version,
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
