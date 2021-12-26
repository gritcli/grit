package deps

import (
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/config"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
	) ([]source.Source, error) {
		var sources []source.Source

		for _, scfg := range cfg.Sources {
			src, err := source.New(scfg)
			if err != nil {
				return nil, err
			}

			sources = append(sources, src)
		}

		return sources, nil
	})
}
