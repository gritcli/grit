package source

import (
	"github.com/gritcli/grit/cmd/gritd/internal/source/githubdriver"
	"github.com/gritcli/grit/internal/config"
)

// Source is an interface for a repository source.
type Source interface {
	// Name returns a short, human-readable identifier of the repository source.
	Name() string

	// Description returns a brief description of the repository source.
	Description() string

	// Close frees any resources allocated for this source.
	Close() error
}

// New returns a new source from a source configuration.
func New(cfg config.Source) (Source, error) {
	var f factory
	cfg.AcceptVisitor(&f)
	return f.Source, f.Error
}

// factory constructs driver-specific sources from source configurations.
type factory struct {
	Source Source
	Error  error
}

func (f *factory) VisitGitHubSource(s config.Source, cfg config.GitHubConfig) {
	f.Source, f.Error = githubdriver.NewSource(s.Name, cfg)
}
