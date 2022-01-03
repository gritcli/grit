package config

import (
	"errors"
	"reflect"

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
	DefaultConfig.Daemon.Socket, _ = homedir.Expand(DefaultConfig.Daemon.Socket)
}

// Config contains an entire Grit configuration.
type Config struct {
	Daemon  Daemon
	Sources map[string]Source
}

// validate returns an error if the configuration is invalid, it is intended to
// be called after any default values have been populated.
func (c Config) validate() error {
	if err := c.Daemon.validate(); err != nil {
		return err
	}

	for _, src := range c.Sources {
		if err := src.validate(); err != nil {
			return err
		}
	}

	return nil
}

// Daemon holds the configuration for the Grit daemon.
type Daemon struct {
	// Socket is the path of the Unix socket used for communication between
	// the Grit CLI and the Grit daemon.
	Socket string `hcl:"socket,optional"`
}

// validate returns an error if the configuration is invalid, it is intended to
// be called after any default values have been populated.
func (d Daemon) validate() error {
	if d.Socket == "" {
		return errors.New("daemon socket must not be empty")
	}

	return nil
}

// Source represents a repository source defined in the configuration.
type Source struct {
	// Name is a short identifier for the source. Each source in the
	// configuration has a unique name.
	Name string

	// Config contains implementation-specific configuration for this source.
	Config SourceConfig
}

// AcceptVisitor calls the appropriate method on v.
func (s Source) AcceptVisitor(v SourceVisitor) {
	s.Config.acceptVisitor(s, v)
}

// validate returns an error if the configuration is invalid, it is intended to
// be called after any default values have been populated.
func (s Source) validate() error {
	if s.Name == "" {
		return errors.New("source name must not be empty")
	}

	return s.Config.validate()
}

// SourceVisitor dispatches Source values to implementation-specific logic.
type SourceVisitor interface {
	VisitGitHubSource(s Source, cfg GitHubConfig)
}

// SourceConfig is an interface for implementation-specific source
// configuration.
type SourceConfig interface {
	// acceptVisitor calls the appropriate method on v.
	acceptVisitor(s Source, v SourceVisitor)

	// withDefaults returns a copy of the configuration with any missing values
	// replaced by their defaults.
	withDefaults() SourceConfig

	// validate returns an error if the configuration is invalid, it is intended
	// to be called after any default values have been populated.
	validate() error
}

// sourceConfigTypes is a map of a source implementation name to the type of its
// SourceConfig implementation.
var sourceConfigTypes = map[string]reflect.Type{}

// registerSourceImpl registers a source implementation, allowing its
// configuration to be parsed.
func registerSourceType(
	name string,
	configType SourceConfig,
	defaultSources ...Source,
) {
	if _, ok := sourceConfigTypes[name]; ok {
		panic("source name already registered")
	}

	sourceConfigTypes[name] = reflect.TypeOf(configType)

	for _, s := range defaultSources {
		DefaultConfig.Sources[s.Name] = s
	}
}
