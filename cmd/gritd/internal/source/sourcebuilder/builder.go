package sourcebuilder

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/cmd/gritd/internal/source/github"
	"github.com/gritcli/grit/internal/config"
)

// Builder builds source.Source values from Grit configuration.
type Builder struct {
	// Logger is the target for log messages from source drivers.
	Logger logging.Logger
}

// FromConfig returns the list of sources defined in cfg.
func (b *Builder) FromConfig(cfg config.Config) []source.Source {
	var sources []source.Source

	for _, cfg := range cfg.Sources {
		src := b.FromSourceConfig(cfg)
		sources = append(sources, src)
	}

	return sources
}

// FromSourceConfig returns the Source defined in cfg.
func (b *Builder) FromSourceConfig(cfg config.Source) source.Source {
	f := &driverFactory{
		Logger: logging.Prefix(
			b.Logger,
			"source[%s]: ",
			cfg.Name,
		),
	}

	cfg.AcceptVisitor(f)

	return source.Source{
		Name:        cfg.Name,
		Description: cfg.DriverConfig.String(),
		Driver:      f.Driver,
	}
}

// driverFactory is a config.SourceVisitor implementation that constructs
// drivers based on source configurations.
type driverFactory struct {
	Logger logging.Logger
	Driver source.Driver
}

func (f *driverFactory) VisitGitHubSource(s config.Source, cfg config.GitHub) {
	d, err := github.NewDriver(cfg, f.Logger)
	if err != nil {
		panic(err)
	}

	f.Driver = d
}
