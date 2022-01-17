package config

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// mergeClonesDefaultsBlock merges b into cfg.
func mergeClonesDefaultsBlock(
	cfg *unresolvedConfig,
	filename string,
	clones clonesSchema,
) error {
	if cfg.ClonesDefaults.File != "" {
		return fmt.Errorf(
			"%s: a 'clones' defaults block is already defined in %s",
			filename,
			cfg.ClonesDefaults.File,
		)
	}

	cfg.ClonesDefaults.File = filename
	cfg.ClonesDefaults.Block = clones

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

// normalizeSourceSpecificClonesBlock normalizes a clonesBlock within a source
// configuration.
func normalizeSourceSpecificClonesBlock(cfg unresolvedConfig, s *unresolvedSource) error {
	if s.Block.ClonesBlock == nil {
		s.Block.ClonesBlock = &clonesSchema{}
	}

	if s.Block.ClonesBlock.Dir == "" {
		s.Block.ClonesBlock.Dir = filepath.Join(
			cfg.ClonesDefaults.Block.Dir,
			s.Block.Name,
		)
	} else {
		// We make sure to only normalize the private key path against s.File if
		// it actually came from the source config (not inherited from the
		// git defaults block).
		if err := normalizePath(s.File, &s.Block.ClonesBlock.Dir); err != nil {
			return err
		}
	}

	return nil
}
