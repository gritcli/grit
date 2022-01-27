package config

import (
	"sort"

	"github.com/gritcli/grit/driver/vcsdriver"
)

// RegisterVCSDriver adds a VCS driver to the registry.
func (r *DriverRegistry) RegisterVCSDriver(alias string, reg vcsdriver.Registration) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.vcsByAlias[alias]; ok {
		panic("alias is already in use")
	}

	if r.vcsByAlias == nil {
		r.vcsByAlias = map[string]vcsdriver.Registration{}
	}

	r.vcsByAlias[alias] = reg
}

// VCSDriverByAlias returns the VCS driver with the given alias.
func (r *DriverRegistry) VCSDriverByAlias(alias string) (vcsdriver.Registration, bool) {
	r.m.RLock()
	reg, ok := r.vcsByAlias[alias]
	r.m.RUnlock()

	if ok {
		return reg, true
	}

	if r.Parent != nil {
		return r.Parent.VCSDriverByAlias(alias)
	}

	return vcsdriver.Registration{}, false
}

// VCSDrivers returns all of the registered vcs drivers.
func (r *DriverRegistry) VCSDrivers() map[string]vcsdriver.Registration {
	drivers := map[string]vcsdriver.Registration{}

	populate := func(r *DriverRegistry) {
		r.m.RLock()
		defer r.m.RUnlock()

		for alias, reg := range r.vcsByAlias {
			drivers[alias] = reg
		}
	}

	if r.Parent != nil {
		populate(r.Parent)
	}

	populate(r)

	return drivers
}

// VCSDriverAliases returns a sorted slice containing the aliases of all
// registered VCS drivers.
func (r *DriverRegistry) VCSDriverAliases() []string {
	drivers := r.VCSDrivers()
	aliases := make([]string, 0, len(drivers))

	for alias := range drivers {
		aliases = append(aliases, alias)
	}

	sort.Strings(aliases)

	return aliases
}
