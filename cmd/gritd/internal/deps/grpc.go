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
	) api.SourceAPIServer {
		return &apiserver.SourceAPIServer{
			Config: cfg,
		}
	})

	Container.Provide(func(
		source api.SourceAPIServer,
	) *grpc.Server {
		s := grpc.NewServer()

		api.RegisterSourceAPIServer(s, source)

		return s
	})

}
