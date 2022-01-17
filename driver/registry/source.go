package registry

import (
	"sort"

	"github.com/gritcli/grit/driver/sourcedriver"
)

// RegisterSourceDriver adds a source driver to the registry.
func (r *Registry) RegisterSourceDriver(alias string, reg sourcedriver.Registration) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.sourceByAlias[alias]; ok {
		panic("alias is already in use")
	}

	if r.sourceByAlias == nil {
		r.sourceByAlias = map[string]sourcedriver.Registration{}
	}

	r.sourceByAlias[alias] = reg
}

// SourceDriverByAlias returns the source driver with the given alias.
func (r *Registry) SourceDriverByAlias(alias string) (sourcedriver.Registration, bool) {
	r.m.RLock()
	reg, ok := r.sourceByAlias[alias]
	r.m.RUnlock()

	if ok {
		return reg, true
	}

	if r.Parent != nil {
		return r.Parent.SourceDriverByAlias(alias)
	}

	return sourcedriver.Registration{}, false
}

// SourceDrivers returns all of the registered source drivers.
func (r *Registry) SourceDrivers() map[string]sourcedriver.Registration {
	drivers := map[string]sourcedriver.Registration{}

	populate := func(r *Registry) {
		r.m.RLock()
		defer r.m.RUnlock()

		for alias, reg := range r.sourceByAlias {
			drivers[alias] = reg
		}
	}

	if r.Parent != nil {
		populate(r.Parent)
	}

	populate(r)

	return drivers
}

// SourceDriverAliases returns a sorted slice containing the aliases of all
// registered source drivers.
func (r *Registry) SourceDriverAliases() []string {
	drivers := r.SourceDrivers()
	aliases := make([]string, 0, len(drivers))

	for alias := range drivers {
		aliases = append(aliases, alias)
	}

	sort.Strings(aliases)

	return aliases
}
