package config

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// Clones is the configuration that controls how Grit stores local repository
// clones.
type Clones struct {
	// Dir is the path to the directory in which local clones are kept.
	Dir string
}

// clonesBlock is the HCL schema for a "clones" block
type clonesBlock struct {
	Dir string `hcl:"dir,optional"`
}

// mergeClonesDefaultsBlock merges b into cfg.
func mergeClonesDefaultsBlock(cfg *unresolvedConfig, filename string, b clonesBlock) error {
	if cfg.ClonesDefaults.File != "" {
		return fmt.Errorf(
			"%s: a 'clones' defaults block is already defined in %s",
			filename,
			cfg.ClonesDefaults.File,
		)
	}

	cfg.ClonesDefaults.File = filename
	cfg.ClonesDefaults.Block = b

	return nil
}

// normalizeClonesDefaultsBlock normalizes cfg.ClonesDefaults.Block and
// populates it with default values.
func normalizeClonesDefaultsBlock(cfg *unresolvedConfig) error {
	if cfg.ClonesDefaults.Block.Dir != "" {
		return normalizePath(cfg.ClonesDefaults.File, &cfg.ClonesDefaults.Block.Dir)
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("unable to determine the current user's home directory")
	}

	cfg.ClonesDefaults.Block.Dir = filepath.Join(homeDir, "grit")

	return nil
}

// assembleClonesBlock converts b into its configuration representation.
func assembleClonesBlock(b clonesBlock) Clones {
	return Clones(b)
}
