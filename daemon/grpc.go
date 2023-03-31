package daemon

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/source"
	"google.golang.org/grpc"
)

func init() {
	imbue.With4(
		container,
		func(
			ctx imbue.Context,
			sources source.List,
			c *source.Cloner,
			s *source.Suggester,
			logger logging.Logger,
		) (api.APIServer, error) {
			return &apiserver.Server{
				SourceList: sources,
				Cloner:     c,
				Suggester:  s,
				Logger:     logging.Prefix(logger, "api: "),
			}, nil
		},
	)

	imbue.With1(
		container,
		func(
			ctx imbue.Context,
			s api.APIServer,
		) (*grpc.Server, error) {
			g := grpc.NewServer()
			api.RegisterAPIServer(g, s)
			return g, nil
		},
	)
}
