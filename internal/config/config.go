package config

import "github.com/mitchellh/go-homedir"

// DefaultDirectory is the default directory to search for Grit configuration
// files.
const DefaultDirectory = "~/.config/grit"

// Config contains an entire Grit configuration.
type Config struct {
	Daemon  Daemon
	Sources map[string]Source
}

// Daemon holds the configuration for the Grit daemon.
type Daemon struct {
	// Socket is the path of the Unix socket used for communication between
	// the Grit CLI and the Grit daemon.
	Socket string `hcl:"socket,optional"`
}

// Source represents a repository source specified in the configuration.
type Source struct {
	// Name is a short identifier for the source. Each source in the
	// configuration has a unique name.
	Name string

	// Driver is the driver used by this source.
	Driver SourceDriver

	// Config contains driver-specific configuration for this source.
	Config DriverConfig
}

// DefaultConfig is the default Grit configuration.
var DefaultConfig = Config{
	Daemon: Daemon{
		Socket: "~/grit/daemon.sock",
	},
	Sources: map[string]Source{},
}

// Normalize the paths in the default configuration.
func init() {
	DefaultConfig.Daemon.Socket, _ = homedir.Expand(DefaultConfig.Daemon.Socket)
}
