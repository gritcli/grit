package source

import (
	"net/url"
	"path"
	"strings"

	"github.com/dogmatiq/dyad"
	"github.com/gritcli/grit/daemon/internal/config"
)

// List is a collection of sources.
type List []Source

// NewList returns a new List from the given source configurations.
func NewList(baseURL *url.URL, sources []config.Source) List {
	var list List

	for _, cfg := range sources {
		if !cfg.Enabled {
			continue
		}

		u := dyad.Clone(baseURL)
		u.Path = path.Join(
			baseURL.Path,
			"source",
			cfg.Name,
		)

		list = append(
			list,
			Source{
				Name:         cfg.Name,
				Description:  cfg.Driver.DescribeSourceConfig(),
				BaseCloneDir: cfg.Clones.Dir,
				BaseURL:      u,
				Driver:       cfg.Driver.NewSource(),
			},
		)
	}

	return list
}

// ByName returns the source with the given name.
//
// Source names are case insensitive.
func (l List) ByName(n string) (Source, bool) {
	for _, src := range l {
		if strings.EqualFold(src.Name, n) {
			return src, true
		}
	}

	return Source{}, false
}
