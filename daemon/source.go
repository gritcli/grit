package daemon

import (
	"net/url"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/daemon/internal/config"
	"github.com/gritcli/grit/daemon/internal/logs"
	"github.com/gritcli/grit/daemon/internal/source"
)

func init() {
	imbue.With2(
		catalog,
		func(
			ctx imbue.Context,
			cfg config.Config,
			baseURL imbue.ByName[httpBaseURL, *url.URL],
		) (source.List, error) {
			return source.NewList(
				baseURL.Value(),
				cfg.Sources,
			), nil
		},
	)

	imbue.With2(
		catalog,
		func(
			ctx imbue.Context,
			sources source.List,
			log logs.Log,
		) (*source.Cloner, error) {
			return &source.Cloner{
				Sources: sources,
				Log:     log,
			}, nil
		},
	)

	imbue.With2(
		catalog,
		func(
			ctx imbue.Context,
			sources source.List,
			log logs.Log,
		) (*source.Suggester, error) {
			return &source.Suggester{
				Sources: sources,
				Log:     log,
			}, nil
		},
	)
}
