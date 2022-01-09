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

// sourceFactory constructs source.Source instances from config.Source values.
type sourceFactory struct {
	Logger logging.Logger
	Source source.Source
	Error  error
}

func (f *sourceFactory) VisitGitHubSource(s config.Source, cfg config.GitHub) {
	d, err := github.NewDriver(
		cfg,
		logging.Prefix(f.Logger, "source[%s]: ", s.Name),
	)

	if err != nil {
		f.Error = err
		return
	}

	f.Source = source.Source{
		Name:   s.Name,
		Driver: d,
	}
}
