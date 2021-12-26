package deps

import (
	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/api"
	"google.golang.org/grpc"
)

func init() {
	Container.Provide(func(
		sources []source.Source,
	) api.APIServer {
		return &apiserver.Server{
			Sources: sources,
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
