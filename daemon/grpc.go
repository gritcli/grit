package daemon

import (
	"os"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/api/daemonapi"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/apiserver/daemonapiserver"
	"github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/logs"
	"google.golang.org/grpc"
)

type version imbue.Name[string]

func init() {
	imbue.With0(
		catalog,
		func(
			ctx imbue.Context,
		) (*grpc.Server, error) {
			return grpc.NewServer(), nil
		},
	)

	imbue.Decorate4(
		catalog,
		func(
			ctx imbue.Context,
			svr *grpc.Server,
			sources source.List,
			c *source.Cloner,
			s *source.Suggester,
			log logs.Log,
		) (*grpc.Server, error) {
			api.RegisterAPIServer(
				svr,
				&apiserver.Server{
					SourceList: sources,
					Cloner:     c,
					Suggester:  s,
					Log:        log.WithPrefix("api: "),
				},
			)
			return svr, nil
		},
	)

	imbue.Decorate1(
		catalog,
		func(
			ctx imbue.Context,
			svr *grpc.Server,
			ver imbue.ByName[version, string],
		) (*grpc.Server, error) {
			daemonapi.RegisterAPIServer(
				svr,
				&daemonapiserver.Server{
					Version: ver.Value(),
					PID:     os.Getpid(),
				},
			)
			return svr, nil
		},
	)
}
