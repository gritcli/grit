package deps

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/gritcli/grit/internal/daemon/internal/apiserver"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"google.golang.org/grpc"
)

func init() {
	Container.Provide(func(
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

	Container.Provide(func(
		s api.APIServer,
	) *grpc.Server {
		g := grpc.NewServer()
		api.RegisterAPIServer(g, s)

		return g
	})
}
