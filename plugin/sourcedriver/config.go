package sourcedriver

// ConfigSchema is an interface for parsing driver-specific configuration within
// a "source" block in a Grit configuration file.
//
// It must be implemented as a pointer-to-struct that uses the field tags
// described by https://pkg.go.dev/github.com/hashicorp/hcl/v2/gohcl, thus it
// defines the HCL schema that is allowed when configuring sources that use this
// driver.
//
// When Grit parses a "source" block within a configuration, any unrecognized
// attributes or blocks within that "source" block are parsed into this schema.
type ConfigSchema interface {
	// Resolve validates and normalizes the configuration to its runtime
	// representation.
	Resolve(ConfigResolutionContext) (Config, error)
}

// ConfigResolutionContext provides operations used to resolve a ConfigSchema to
// a Config.
type ConfigResolutionContext interface {
	// ResolveVCSConfig resolves a driver-specific configuration block for one
	// of the supported version control systems.
	ResolveVCSConfig(in, out interface{}) error
}

// Config is an interface for driver-specific configuration options for a
// repository source.
type Config interface {
	// NewDriver constructs a new driver that uses this configuration.
	NewDriver() Driver

	// String returns a human-readable description of the configuration.
	String() string
}
