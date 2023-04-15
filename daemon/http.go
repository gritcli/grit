package daemon

import (
	"net"
	"net/http"

	"github.com/dogmatiq/imbue"
)

type httpServer imbue.Group

func init() {
	imbue.With0Grouped[httpServer](
		catalog,
		func(
			ctx imbue.Context,
		) (net.Listener, error) {
			lis, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return nil, err
			}
			ctx.Defer(lis.Close)

			return lis, nil
		},
	)

	imbue.With0Grouped[httpServer](
		catalog,
		func(
			ctx imbue.Context,
		) (*http.ServeMux, error) {
			return http.NewServeMux(), nil
		},
	)
}
