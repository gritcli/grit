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
	// Normalize validates the configuration as parsed by this schema and
	// returns a normalized Config.
	//
	// The implementation must call ctx.ReadVCSConfig() for each VCS driver
	// that is supported, even if they are not currently in use.
	Normalize(ctx ConfigNormalizeContext) (Config, error)
}

// ConfigNormalizeContext provides operations used to normalize a
// ConfigSchema.
type ConfigNormalizeContext interface {
	// NormalizePath normalizes a filesystem encountered within the
	// configuration.
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
