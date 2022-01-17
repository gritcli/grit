package config

// DefaultDirectory is the default directory to search for Grit configuration
// files.
const DefaultDirectory = "~/.config/grit"

// Config contains an entire Grit configuration.
type Config struct {
	// Daemon is the configuration of the Grit daemon.
	Daemon Daemon

	// Sources is the set of repository sources from which repositories can be
	// cloned.
	Sources []Source
}

// configFile is HCL schema for a single configuration file.
type configFile struct {
	DaemonBlock         *daemonBlock  `hcl:"daemon,block"`
	ClonesDefaultsBlock *clonesBlock  `hcl:"clones,block"`
	VCSDefaultsBlocks   []vcsBlock    `hcl:"vcs,block"`
	SourceBlocks        []sourceBlock `hcl:"source,block"`
}
