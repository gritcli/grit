package config

// DefaultDirectory is the default directory to search for Grit configuration
// files.
const DefaultDirectory = "~/.config/grit"

// Config contains an entire Grit configuration.
type Config struct {
	// Daemon is the configuration of the Grit daemon.
	Daemon Daemon

	// GitDefaults is the configuration that controls how Grit uses Git across all
	// sources. Repository sources that use Git may allow these settings to be
	// overridden.
	GitDefaults Git

	// Sources is the set of repository sources from which repositories can be
	// cloned.
	Sources map[string]Source
}

// configFile is HCL schema for a single configuration file.
type configFile struct {
	DaemonBlock      *daemonBlock  `hcl:"daemon,block"`
	GitDefaultsBlock *gitBlock     `hcl:"git,block"`
	SourceBlocks     []sourceBlock `hcl:"source,block"`
}
