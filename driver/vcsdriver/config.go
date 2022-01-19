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
	// NormalizeGlobals validates the global configuration as parsed by this
	// schema, and returns a normalized Config.
	//
	// The "global" VCS configuration is a "vcs" block that appears at the
	// "top-level" of a configuration file. Such global configuration informs
	// the source-specific VCS configuration.
	//
	// If there is no "vcs" block for this driver at the top-level of the Grit
	// configuration, the method is called on a zero-value ConfigSchema.
	NormalizeGlobals(ctx ConfigNormalizeContext) (Config, error)

	// NormalizeSourceSpecific validates the configuration as parsed by this
	// schema within a "source" block and returns a normalized Config.
	//
	// g is the global configuration for this VCS, as returned by
	// NormalizeGlobals().
	NormalizeSourceSpecific(ctx ConfigNormalizeContext, g Config) (Config, error)
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
//
// The underlying implementation must not be used by more than one driver.
type Config interface {
	// DescribeVCSConfig returns a human-readable description of the
	// configuration.
	DescribeVCSConfig() string
}
