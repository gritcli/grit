package registry

import (
	"sort"

	"github.com/gritcli/grit/plugin/sourcedriver"
)

// RegisterSourceDriver adds a source driver to the registry.
func (r *Registry) RegisterSourceDriver(alias string, d sourcedriver.Registration) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.sourceByAlias[alias]; ok {
		panic("alias is already in use")
	}

	if r.sourceByAlias == nil {
		r.sourceByAlias = map[string]sourcedriver.Registration{}
	}

	r.sourceByAlias[alias] = d
}

// SourceDriverByAlias returns the source driver with the given alias.
func (r *Registry) SourceDriverByAlias(alias string) (sourcedriver.Registration, bool) {
	r.m.RLock()
	d, ok := r.sourceByAlias[alias]
	r.m.RUnlock()

	if ok {
		return d, true
	}

	if r.Parent != nil {
		return r.Parent.SourceDriverByAlias(alias)
	}

	return sourcedriver.Registration{}, false
}

// SourceDriverAliases returns the aliases of all registered drivers.
func (r *Registry) SourceDriverAliases() []string {
	uniq := map[string]struct{}{}

	populate := func(r *Registry) {
		r.m.RLock()
		defer r.m.RUnlock()

		for alias := range r.sourceByAlias {
			uniq[alias] = struct{}{}
		}
	}

	populate(r)

	if r.Parent != nil {
		populate(r.Parent)
	}

	aliases := make([]string, 0, len(uniq))
	for alias := range uniq {
		aliases = append(aliases, alias)
	}

	sort.Strings(aliases)

	return aliases
}
