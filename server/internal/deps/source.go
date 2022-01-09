package deps

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/common/config"
	"github.com/gritcli/grit/server/internal/source"
	"github.com/gritcli/grit/server/internal/source/sourcebuilder"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
		logger logging.Logger,
	) []source.Source {
		builder := &sourcebuilder.Builder{
			Logger: logger,
		}

		return builder.FromConfig(cfg)
	})
}
