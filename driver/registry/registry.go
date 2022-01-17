package registry

import (
	"sync"

	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
)

// BuiltIns is the registry of official drivers that ship with Grit.
var BuiltIns Registry

// Registry is a collection of driver implementations.
type Registry struct {
	Parent *Registry

	m             sync.RWMutex
	sourceByAlias map[string]sourcedriver.Registration
	vcsByAlias    map[string]vcsdriver.Registration
}
