package deps

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/cmd/gritd/internal/source/github"
	"github.com/gritcli/grit/internal/config"
)

func init() {
	Container.Provide(func(
		cfg config.Config,
		logger logging.Logger,
	) ([]source.Source, error) {
		var sources []source.Source

		factory := sourceFactory{
			Logger: logger,
		}

		for _, c := range cfg.Sources {
			c.AcceptVisitor(&factory)

			if factory.Error != nil {
				return nil, factory.Error
			}

			sources = append(sources, factory.Source)
		}

		return sources, nil
	})
}

// sourceFactory constructs driver-specific sources from source configurations.
type sourceFactory struct {
	Logger logging.Logger
	Source source.Source
	Error  error
}

func (f *sourceFactory) VisitGitHubSource(s config.Source, cfg config.GitHubConfig) {
	f.Source, f.Error = github.NewSource(
		s.Name,
		cfg,
		logging.Prefix(f.Logger, "source[%s]: ", s.Name),
	)
}
