package config

import (
	"errors"
	"fmt"
)

// Git is the configuration that controls how Grit uses Git.
type Git struct {
	// PrivateKey is the path to the private SSH key used to authenticate when
	// using the SSH transport. If it is empty, the system's SSH agent is
	// queried to determine which identity to use.
	PrivateKey string

	// Passphrase is the passphrase used to encrypt the private key, if any.
	Passphrase string

	// PreferHTTP indicates that Grit should prefer the HTTP transport. By
	// default SSH is preferred.
	PreferHTTP bool
}

// gitBlock is the HCL schema for a "git" block
type gitBlock struct {
	PrivateKey string `hcl:"private_key,optional"`
	Passphrase string `hcl:"passphrase,optional"`
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

	if err := validateGitBlock(b); err != nil {
		return fmt.Errorf(
			"%s: the global 'git' block is invalid: %w",
			filename,
			err,
		)
	}

	return nil
}

// validateGitBlock validates the contents of b.
func validateGitBlock(b gitBlock) error {
	if b.PrivateKey == "" && b.Passphrase != "" {
		return errors.New("passphrase present without specifying a private key file")
	}

	return nil
}

// normalizeGlobalGitBlock normalizes cfg.GlobalGit.Block and populates it with
// default values.
func normalizeGlobalGitBlock(cfg *unresolvedConfig) error {
	return normalizePath(cfg.GlobalGit.File, &cfg.GlobalGit.Block.PrivateKey)
}

// normalizeSourceSpecificGitBlock normalizes a gitBlock within a source
// configuration.
func normalizeSourceSpecificGitBlock(cfg unresolvedConfig, s unresolvedSource, p **gitBlock) error {
	if *p == nil {
		*p = &gitBlock{}
	}

	b := *p

	if b.PrivateKey == "" {
		b.PrivateKey = cfg.GlobalGit.Block.PrivateKey

		if b.Passphrase == "" {
			b.Passphrase = cfg.GlobalGit.Block.Passphrase
		}
	} else {
		// We make sure to only normalize the private key path against s.File if
		// it actually came from the source config (not inherited from the
		// global git block).
		if err := normalizePath(s.File, &b.PrivateKey); err != nil {
			return err
		}
	}

	if b.PreferHTTP == nil {
		b.PreferHTTP = cfg.GlobalGit.Block.PreferHTTP
	}

	if err := validateGitBlock(*b); err != nil {
		return fmt.Errorf(
			"the 'git' block is invalid: %w",
			err,
		)
	}

	return nil
}

// assembleGitBlock converts b into its configuration representation.
func assembleGitBlock(b gitBlock) Git {
	cfg := Git{
		PrivateKey: b.PrivateKey,
		Passphrase: b.Passphrase,
	}

	if b.PreferHTTP != nil {
		cfg.PreferHTTP = *b.PreferHTTP
	}

	return cfg
}
