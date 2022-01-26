package vcsdriver

import "github.com/hashicorp/hcl/v2"

// ConfigLoader is an interface for loading driver-specific VCS configuration.
type ConfigLoader interface {
	// Defaults returns the default configuration for this driver.
	Defaults(ctx ConfigContext) (Config, error)

	// UnmarshalAndMerge unmarshals the contents of a "vcs" block and returns
	// the result of merging it with an existing configuration.
	//
	// c is the existing configuration, b is the body of the "vcs" block. c must
	// not be modified.
	UnmarshalAndMerge(ctx ConfigContext, c Config, b hcl.Body) (Config, error)
}

// ConfigContext provides operations used when loading VCS configuration.
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
}

// Config is a driver-specific VCS configuration.
//
// The underlying implementation must not be used by more than one VCS driver.
type Config interface {
	// DescribeVCSConfig returns a human-readable description of the
	// configuration.
	DescribeVCSConfig() string
}
