package daemon

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/daemon/internal/source"
)

func init() {
	container.Provide(func(
		cfg config.Config,
	) source.List {
		return source.NewList(cfg.Sources)
	})

	container.Provide(func(
		sources source.List,
		logger logging.Logger,
	) *source.Cloner {
		return &source.Cloner{
			Sources: sources,
			Logger:  logger,
		}
	})

	container.Provide(func(
		sources source.List,
		logger logging.Logger,
	) *source.Suggester {
		return &source.Suggester{
			Sources: sources,
			Logger:  logger,
		}
	})
}
