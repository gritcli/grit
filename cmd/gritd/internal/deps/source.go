package deps

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/cmd/gritd/internal/source/sourcebuilder"
	"github.com/gritcli/grit/common/config"
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
