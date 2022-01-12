package deps

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/daemon/internal/config"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/daemon/internal/source/sourcebuilder"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
	) source.List {
		builder := &sourcebuilder.Builder{}
		return builder.FromConfig(cfg)
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
