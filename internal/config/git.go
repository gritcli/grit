package config

import "fmt"

// Git is the configuration that controls how Grit uses Git.
type Git struct {
	// PrivateKey is the path to the private SSH key used to authenticate when
	// using the SSH transport. If it is empty, the system's SSH agent is
	// queried to determine which identity to use.
	PrivateKey string

	// PreferHTTP indicates that Grit should prefer the HTTP transport. By
	// default SSH is preferred.
	PreferHTTP bool
}

// gitBlock is the HCL schema for a "git" block
type gitBlock struct {
	PrivateKey string `hcl:"private_key,optional"`
	PreferHTTP *bool  `hcl:"prefer_http"` // pointer allows detection of absence vs explicit false
}

// mergeGlobalGitBlock merges b into cfg.
func mergeGlobalGitBlock(cfg *unresolvedConfig, filename string, b gitBlock) error {
	if cfg.GlobalGit.File != "" {
		return fmt.Errorf(
			"%s: a global 'git' block is already defined in %s",
			filename,
			cfg.GlobalGit.File,
		)
	}

	cfg.GlobalGit.File = filename
	cfg.GlobalGit.Block = b

	return nil
}

// normalizeGlobalGitBlock normalizes cfg.GlobalGit.Block and populates it with
// default values.
func normalizeGlobalGitBlock(cfg *unresolvedConfig) error {
	return normalizePath(cfg.GlobalGit.File, &cfg.GlobalGit.Block.PrivateKey)
}

// assembleGitBlock converts b into its configuration representation.
func assembleGitBlock(b gitBlock) Git {
	cfg := Git{
		PrivateKey: b.PrivateKey,
	}

	if b.PreferHTTP != nil {
		cfg.PreferHTTP = *b.PreferHTTP
	}

	return cfg
}
