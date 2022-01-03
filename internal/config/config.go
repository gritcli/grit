package config

import (
	"github.com/mitchellh/go-homedir"
)

// DefaultDirectory is the default directory to search for Grit configuration
// files.
const DefaultDirectory = "~/.config/grit"

// DefaultConfig is the default Grit configuration.
var DefaultConfig = Config{
	Daemon: Daemon{
		Socket: "~/grit/daemon.sock",
	},
	Sources: map[string]Source{},
}

// Normalize the paths in the default configuration.
func init() {
	var err error
	DefaultConfig.Daemon.Socket, err = homedir.Expand(DefaultConfig.Daemon.Socket)
	if err != nil {
		panic(err)
	}
}

// Config contains an entire Grit configuration.
type Config struct {
	// Daemon is the configuration of the Grit daemon.
	Daemon Daemon

	// Sources is the set of repository sources from which repositories can be
	// cloned.
	Sources map[string]Source
}

// Daemon holds the configuration for the Grit daemon.
type Daemon struct {
	// Socket is the path of the Unix socket used for communication between
	// the Grit CLI and the Grit daemon.
	Socket string
}

// Source represents a repository source defined in the configuration.
type Source struct {
	// Name is a short identifier for the source. Each source in the
	// configuration has a unique name.
	Name string

	// Enabled is true if the source is enabled. Disabled sources are not used
	// when searching for repositories to be cloned.
	Enabled bool

	// Config contains implementation-specific configuration for this source.
	Config SourceConfig
}

// AcceptVisitor calls the appropriate method on v.
func (s Source) AcceptVisitor(v SourceVisitor) {
	s.Config.acceptVisitor(s, v)
}

// SourceConfig is an interface for implementation-specific configuration
// options for a repository source.
type SourceConfig interface {
	// acceptVisitor calls the appropriate method on v.
	acceptVisitor(s Source, v SourceVisitor)
}

// SourceVisitor dispatches Source values to implementation-specific logic.
type SourceVisitor interface {
	VisitGitHubSource(s Source, cfg GitHubConfig)
}
