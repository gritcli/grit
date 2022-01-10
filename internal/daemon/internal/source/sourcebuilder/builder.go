package sourcebuilder

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/common/config"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/daemon/internal/source/internal/github"
)

// Builder builds source.Source values from Grit configuration.
type Builder struct {
	// Logger is the target for log messages from source drivers.
	Logger logging.Logger
}

// FromConfig returns the list of enabled sources defined in cfg.
func (b *Builder) FromConfig(cfg config.Config) []source.Source {
	var sources []source.Source

	for _, cfg := range cfg.Sources {
		if cfg.Enabled {
			src := b.FromSourceConfig(cfg)
			sources = append(sources, src)
		}
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
		Description: cfg.Driver.String(),
		CloneDir:    cfg.Clones.Dir,
		Driver:      f.Driver,
	}
}

// driverFactory is a config.SourceVisitor implementation that constructs
// drivers based on source configurations.
type driverFactory struct {
	Driver source.Driver
	Logger logging.Logger
}

func (f *driverFactory) VisitGitHubSource(s config.Source, cfg config.GitHub) {
	f.Driver = &github.Driver{
		Config: cfg,
		Logger: f.Logger,
	}
}
