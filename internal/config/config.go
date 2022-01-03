package config

import (
	"errors"
	"reflect"
	"regexp"

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
	Daemon  Daemon
	Sources map[string]Source
}

// Daemon holds the configuration for the Grit daemon.
type Daemon struct {
	// Socket is the path of the Unix socket used for communication between
	// the Grit CLI and the Grit daemon.
	Socket string `hcl:"socket,optional"`
}

// normalize validates the configuration and returns a copy with any missing
// values replaced by their defaults.
func (d Daemon) normalize(filename string) (Daemon, error) {
	if d.Socket == "" {
		d.Socket = DefaultConfig.Daemon.Socket
	}

	return d, normalizePath(filename, &d.Socket)
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

var sourceNameRegexp = regexp.MustCompile(`(?i)^[a-z_]+$`)

// normalize validates the configuration and returns a copy with any missing
// values replaced by their defaults.
func (s Source) normalize(filename string) (Source, error) {
	if s.Name == "" {
		return Source{}, errors.New("source name must not be empty")
	}

	if !sourceNameRegexp.MatchString(s.Name) {
		return Source{}, errors.New("source name must contain only alpha-numeric characters and underscores")
	}

	return s, nil
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

	// normalize validates the configuration and returns a copy with any missing
	// values replaced by their defaults.
	normalize(filename string) (SourceConfig, error)
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
