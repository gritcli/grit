package deps

import (
	"github.com/gritcli/grit/internal/common/api"
	"github.com/gritcli/grit/internal/daemon/internal/apiserver"
	"google.golang.org/grpc"
)

func init() {
	Container.Provide(apiserver.New)

	Container.Provide(func(
		s api.APIServer,
	) *grpc.Server {
		g := grpc.NewServer()
		api.RegisterAPIServer(g, s)

		return g
	})

}
