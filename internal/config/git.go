package config

import (
	"fmt"
)

// Git is the configuration that controls how Grit uses Git.
type Git struct {
	// SSHKeyFile is the path to the private SSH key used to authenticate when
	// using the SSH transport.
	//
	// If it is empty, the system's SSH agent is queried to determine which
	// identity to use.
	SSHKeyFile string

	// SSHKeyPassphrase is the passphrase used to encrypt the SSH private key,
	// if any.
	SSHKeyPassphrase string

	// PreferHTTP indicates that Grit should prefer the HTTP transport. By
	// default SSH is preferred.
	PreferHTTP bool
}

// gitBlock is the HCL schema for a "git" block
type gitBlock struct {
	SSHKey     *sshKeyBlock `hcl:"ssh_key,block"`
	PreferHTTP *bool        `hcl:"prefer_http"` // pointer allows detection of absence vs explicit false
}

// sshKeyBlock is the HCL schema for a "private_key" block within a "git"
// block.
type sshKeyBlock struct {
	File       string `hcl:"file"`
	Passphrase string `hcl:"passphrase,optional"`
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
	return nil
}

// normalizeGlobalGitBlock normalizes cfg.GlobalGit.Block and populates it with
// default values.
func normalizeGlobalGitBlock(cfg *unresolvedConfig) error {
	if cfg.GlobalGit.Block.SSHKey != nil {
		return normalizePath(cfg.GlobalGit.File, &cfg.GlobalGit.Block.SSHKey.File)
	}

	return nil
}

// normalizeSourceSpecificGitBlock normalizes a gitBlock within a source
// configuration.
func normalizeSourceSpecificGitBlock(cfg unresolvedConfig, s unresolvedSource, p **gitBlock) error {
	if *p == nil {
		*p = &gitBlock{}
	}

	b := *p

	if b.SSHKey == nil {
		b.SSHKey = cfg.GlobalGit.Block.SSHKey
	} else {
		// We make sure to only normalize the private key path against s.File if
		// it actually came from the source config (not inherited from the
		// global git block).
		if err := normalizePath(s.File, &b.SSHKey.File); err != nil {
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
	cfg := Git{}

	if b.SSHKey != nil {
		cfg.SSHKeyFile = b.SSHKey.File
		cfg.SSHKeyPassphrase = b.SSHKey.Passphrase
	}

	if b.PreferHTTP != nil {
		cfg.PreferHTTP = *b.PreferHTTP
	}

	return cfg
}
