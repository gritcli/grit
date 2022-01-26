package sourcedriver

import "github.com/hashicorp/hcl/v2"

// ImplicitSource represents a source that is provided by this driver without
// explicit configuration.
type ImplicitSource struct {
	// Name is the unique name for the source.
	//
	// If the user defines an explicit source with the same name, this implicit
	// source is ignored.
	Name string

	// Config is the configuration of this source.
	Config Config
}

// ConfigLoader is an interface for loading driver-specific source configuration.
type ConfigLoader interface {
	// Unmarshal unmarshals the contents of a "source" block.
	Unmarshal(ctx ConfigContext, b hcl.Body) (Config, error)

	// ImplicitSources returns the configuration to use for "implicit" sources
	// provided by this driver without explicit configuration.
	ImplicitSources(ctx ConfigContext) ([]ImplicitSource, error)
}

// ConfigContext provides operations used when loading source configuration.
type ConfigContext interface {
	// EvalContext returns the HCL evaluation context to be used when to
	// decoding HCL content.
	EvalContext() *hcl.EvalContext

	// NormalizePath resolves a (potentially relative) filesystem path to an
	// absolute path.
	//
	// If *p begins with a tilde (~), it is resolved relative to the user's home
	// directory.
	//
	// If *p is a relative path, it is resolved to an absolute path relative to
	// the directory containing the configuration file that is currently being
	// parsed.
	//
	// It does nothing if p is nil or *p is empty.
	NormalizePath(p *string) error

	// UnmarshalVCSConfig stores the configuration for a VCS driver in the value
	// pointed to by v.
	//
	// It panics v is nil or not a pointer.
	//
	// driver is the name (not the alias) of the VCS driver that provides the
	// configuration, as specified in its vcsdriver.Registration entry.
	//
	// Multiple drivers may share the same name by using different aliases. In
	// this case, the type of the value pointed to by v the driver that supplies
	// a config of the type pointed to by v is used.
	UnmarshalVCSConfig(driver string, v interface{}) error
}

// Config is an interface for driver-specific configuration options for a
// repository source.
//
// The underlying implementation must not be used by more than one driver.
type Config interface {
	// NewDriver constructs a new driver that uses this configuration.
	NewDriver() Driver

	// DescribeSourceConfig returns a human-readable description of the
	// configuration.
	DescribeSourceConfig() string
}
