package vcsdriver

import "github.com/hashicorp/hcl/v2"

// ConfigNormalizer is an interface for normalizing driver-specific
// configuration within a "vcs" block in a Grit configuration file.
type ConfigNormalizer interface {
	// Defaults returns the default configuration to use for this driver.
	Defaults(nc ConfigNormalizeContext) (Config, error)

	// Merge returns a new Config that is the result of merging an existing
	// Config with the contents of a "vcs" block.
	//
	// c is the existing configuration, b is the body of the "vcs" block. c must
	// not be modified.
	Merge(nc ConfigNormalizeContext, c Config, b hcl.Body) (Config, error)
}

// ConfigNormalizeContext provides operations used to normalize a
// ConfigSchema.
type ConfigNormalizeContext interface {
	// EvalContext returns the HCL evaluation context to be used when to
	// decoding HCL content.
	EvalContext() *hcl.EvalContext

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
