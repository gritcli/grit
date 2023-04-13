package daemon

import (
	"net"
	"net/url"

	"github.com/dogmatiq/imbue"
)

type (
	httpListener imbue.Name[net.Listener]
	httpBaseURL  imbue.Name[*url.URL]
)

func init() {
	imbue.With0Named[httpListener](
		catalog,
		func(
			ctx imbue.Context,
		) (net.Listener, error) {
			lis, err := net.Listen("tcp", "localhost:0")
			if err != nil {
				return nil, err
			}
			ctx.Defer(lis.Close)

			return lis, nil
		},
	)

	imbue.With1Named[httpBaseURL](
		catalog,
		func(
			ctx imbue.Context,
			lis imbue.ByName[httpListener, net.Listener],
		) (*url.URL, error) {
			return &url.URL{
				Scheme: "http",
				Host:   lis.Value().Addr().String(),
			}, nil
		},
	)
}
