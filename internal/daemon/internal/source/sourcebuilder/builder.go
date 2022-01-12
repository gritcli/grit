package sourcebuilder

import (
	"sort"

	"github.com/gritcli/grit/internal/daemon/internal/config"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/daemon/internal/source/internal/github"
)

// Builder builds source.Source values from Grit configuration.
type Builder struct {
}

// FromConfig returns the list of enabled sources defined in cfg.
func (b *Builder) FromConfig(cfg config.Config) source.List {
	var sources source.List

	for _, cfg := range cfg.Sources {
		if cfg.Enabled {
			src := b.FromSourceConfig(cfg)
			sources = append(sources, src)
		}
	}

	sort.Slice(
		sources,
		func(i, j int) bool {
			return sources[i].Name < sources[j].Name
		},
	)

	return sources
}

// FromSourceConfig returns the Source defined in cfg.
func (b *Builder) FromSourceConfig(cfg config.Source) source.Source {
	var f driverFactory

	cfg.AcceptVisitor(&f)

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
}

func (f *driverFactory) VisitGitHubSource(s config.Source, cfg config.GitHub) {
	f.Driver = &github.Driver{
		Config: cfg,
	}
}
