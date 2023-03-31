package daemon

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/daemon/internal/source"
)

func init() {
	imbue.With1(
		container,
		func(
			ctx imbue.Context,
			cfg config.Config,
		) (source.List, error) {
			return source.NewList(cfg.Sources), nil
		},
	)

	imbue.With2(
		container,
		func(
			ctx imbue.Context,
			s source.List,
			l logging.Logger,
		) (*source.Cloner, error) {
			return &source.Cloner{
				Sources: s,
				Logger:  l,
			}, nil
		},
	)

	imbue.With2(
		container,
		func(
			ctx imbue.Context,
			s source.List,
			l logging.Logger,
		) (*source.Suggester, error) {
			return &source.Suggester{
				Sources: s,
				Logger:  l,
			}, nil
		},
	)
}
