package deps

import (
	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/common/api"
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
