package registry

import (
	"sort"

	"github.com/gritcli/grit/driver/vcsdriver"
)

// RegisterVCSDriver adds a VCS driver to the registry.
func (r *Registry) RegisterVCSDriver(alias string, v vcsdriver.Registration) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.vcsByAlias[alias]; ok {
		panic("alias is already in use")
	}

	if r.vcsByAlias == nil {
		r.vcsByAlias = map[string]vcsdriver.Registration{}
	}

	r.vcsByAlias[alias] = v
}

// VCSDriverByAlias returns the VCS driver with the given alias.
func (r *Registry) VCSDriverByAlias(alias string) (vcsdriver.Registration, bool) {
	r.m.RLock()
	d, ok := r.vcsByAlias[alias]
	r.m.RUnlock()

	if ok {
		return d, true
	}

	if r.Parent != nil {
		return r.Parent.VCSDriverByAlias(alias)
	}

	return vcsdriver.Registration{}, false
}

// VCSDriverAliases returns the aliases of all registered drivers.
func (r *Registry) VCSDriverAliases() []string {
	uniq := map[string]struct{}{}

	populate := func(r *Registry) {
		r.m.RLock()
		defer r.m.RUnlock()

		for alias := range r.vcsByAlias {
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
