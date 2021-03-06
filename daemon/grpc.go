package daemon

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/source"
	"google.golang.org/grpc"
)

func init() {
	container.Provide(func(
		sources source.List,
		c *source.Cloner,
		s *source.Suggester,
		logger logging.Logger,
	) api.APIServer {
		return &apiserver.Server{
			SourceList: sources,
			Cloner:     c,
			Suggester:  s,
			Logger:     logging.Prefix(logger, "api: "),
		}
	})

	container.Provide(func(
		s api.APIServer,
	) *grpc.Server {
		g := grpc.NewServer()
		api.RegisterAPIServer(g, s)

		return g
	})
}
