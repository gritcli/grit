package daemon

import (
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/logs"
	"google.golang.org/grpc"
)

func init() {
	imbue.With4(
		catalog,
		func(
			ctx imbue.Context,
			sources source.List,
			c *source.Cloner,
			s *source.Suggester,
			log logs.Log,
		) (api.APIServer, error) {
			return &apiserver.Server{
				SourceList: sources,
				Cloner:     c,
				Suggester:  s,
				Log:        log.WithPrefix("api: "),
			}, nil
		},
	)

	imbue.With1(
		catalog,
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
