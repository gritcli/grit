package deps

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/daemon/internal/config"
	"github.com/gritcli/grit/internal/daemon/internal/source"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
	) source.List {
		return source.NewList(cfg.Sources)
	})

	Container.Provide(func(
		sources source.List,
		logger logging.Logger,
	) *source.Cloner {
		return &source.Cloner{
			Sources: sources,
			Logger:  logger,
		}
	})
}
