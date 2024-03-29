package config

import (
	"path/filepath"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
)

var (
	// DefaultDirectory is the standard location for Grit configuration files.
	DefaultDirectory = filepath.Join("~", ".config", "grit")

	// DefaultClonesDirectory is the default path in which grit stores local
	// clones of remote repositories.
	DefaultClonesDirectory = filepath.Join("~", "grit")
)

// Config contains an entire Grit configuration.
type Config struct {
	// Daemon is the configuration of the Grit daemon.
	Daemon Daemon

	// Sources is the set of repository sources from which repositories can be
	// cloned.
	Sources []Source
}

// Daemon holds the configuration for the Grit daemon.
type Daemon struct {
	// Socket is the path of the Unix socket used for communication between
	// the Grit CLI and the Grit daemon (via gRPC).
	Socket string
}

// Source is the configuration for a source of repositories.
type Source struct {
	// Name is a short identifier for the source. Each source in the
	// configuration has a unique name. Names are case-insensitive.
	Name string

	// Enabled is true if the source is enabled. Disabled sources are not used
	// when searching for repositories to be cloned.
	Enabled bool

	// Clones is the configuration that controls how Grit stores local
	// repository clones for this source.
	Clones Clones

	// Driver contains driver-specific configuration for this source.
	Driver sourcedriver.Config
}

// Clones is the configuration that controls how Grit stores local clones of
// remote repositories.
type Clones struct {
	// Dir is the path to the directory in which local clones are kept.
	Dir string
}
