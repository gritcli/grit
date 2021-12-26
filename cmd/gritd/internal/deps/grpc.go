package deps

import (
	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/config"
	"google.golang.org/grpc"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
	) api.APIServer {
		return &apiserver.Server{
			Config: cfg,
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
