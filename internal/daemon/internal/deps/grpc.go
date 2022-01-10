package deps

import (
	"github.com/gritcli/grit/internal/common/api"
	"github.com/gritcli/grit/internal/daemon/internal/apiserver"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"google.golang.org/grpc"
)

func init() {
	Container.Provide(func(
		sources source.List,
	) api.APIServer {
		return &apiserver.Server{
			SourceList: sources,
		}
	})

	Container.Provide(func(
		s api.APIServer,
	) *grpc.Server {
		g := grpc.NewServer()
		api.RegisterAPIServer(g, s)

		return g
	})

}
