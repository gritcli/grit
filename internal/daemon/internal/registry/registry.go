package registry

import (
	"sync"

	"github.com/gritcli/grit/plugin/driver"
)

// BuiltIns is the registry of official drivers that ship with Grit.
var BuiltIns Registry

// Registry is a collection of driver implementations.
type Registry struct {
	Parent *Registry

	m             sync.RWMutex
	sourceByAlias map[string]driver.Registration
}
