package vcsdriver

// ConfigSchema is an interface for parsing driver-specific configuration within
// a "vcs" block in a Grit configuration file.
//
// It must be implemented as a pointer-to-struct that uses the field tags
// described by https://pkg.go.dev/github.com/hashicorp/hcl/v2/gohcl, thus it
// defines the HCL schema that is allowed when configuring version control
// systems that use this driver.
//
// When Grit parses a "vcs" block within a configuration, any unrecognized
// attributes or blocks within that "vcs" block are parsed into this schema.
type ConfigSchema interface {
	// NormalizeDefaults validates the configuration as parsed by this schema at
	// the "top-level" of a Grit configuration, and returns a normalized Config.
	//
	// VCS configuration that appears at the "top-level" of the configuration is
	// used to set the defaults used by any source driver that uses this VCS.
	//
	// If there is no "vcs" block for this driver at the top-level of the Grit
	// configuration, the method is called on a zero-value ConfigSchema.
	NormalizeDefaults(ctx ConfigNormalizeContext) (Config, error)

	// NormalizeSourceSpecific validates the configuration as parsed by this
	// schema within a "source" block and returns a normalized Config.
	//
	// defaults is the result of a prior call to NormalizeDefaults().
	NormalizeSourceSpecific(ctx ConfigNormalizeContext, defaults Config) (Config, error)
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
}

// Config is an interface for driver-specific configuration options for a VCS.
type Config interface {
	// DescribeVCSConfig returns a human-readable description of the
	// configuration.
	DescribeVCSConfig() string
}
