package config

import (
	"sync"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
)

// DriverRegistry is a collection of driver implementations that are available
// to the configuration.
type DriverRegistry struct {
	Parent *DriverRegistry

	m             sync.RWMutex
	sourceByAlias map[string]sourcedriver.Registration
	vcsByAlias    map[string]vcsdriver.Registration
}
