package config

// DefaultDirectory is the default directory to search for Grit configuration
// files.
const DefaultDirectory = "~/.config/grit"

// Config is the root of a Grit configuration.
type Config struct {
	Daemon Daemon
}

// mainConfig is the structure of the "grit.hcl" file.
type mainConfig struct {
	Daemon *Daemon `hcl:"daemon,block"`
}

// sourceConfig is the structure of a "*.source.hcl" file.
type sourceConfig struct {
}

// Daemon contains configuration for the Grit daemon.
type Daemon struct {
	// Socket is the path of the Unix socket used for communication between
	// the Grit CLI and the Grit daemon.
	Socket string `hcl:"socket"`
}

// DefaultConfig is the default Grit configuration.
var DefaultConfig = Config{
	Daemon: Daemon{
		Socket: "~/grit/grit.sock",
	},
}
