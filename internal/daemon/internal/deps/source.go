package deps

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/common/config"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/daemon/internal/source/sourcebuilder"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
		logger logging.Logger,
	) source.List {
		builder := &sourcebuilder.Builder{
			Logger: logger,
		}

		return builder.FromConfig(cfg)
	})
}
