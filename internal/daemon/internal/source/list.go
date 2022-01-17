package source

import (
	"sort"
	"strings"

	"github.com/gritcli/grit/internal/daemon/internal/config"
)

// List is a collection of sources.
type List []Source

// NewList returns a new List from the given source configurations.
func NewList(sources []config.Source) List {
	var list List

	for _, cfg := range sources {
		if !cfg.Enabled {
			continue
		}

		list = append(
			list,
			Source{
				Name:        cfg.Name,
				Description: cfg.Driver.String(),
				CloneDir:    cfg.Clones.Dir,
				Driver:      cfg.Driver.NewDriver(),
			},
		)
	}

	sort.Slice(
		list,
		func(i, j int) bool {
			return list[i].Name < list[j].Name
		},
	)

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
