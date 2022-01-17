package config

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// mergeRootClones merges s into the configuration.
//
// s is a "clones" block at the top-level of the current file (that is, not
// within a "source" block).
func (r *resolver) mergeRootClones(s clonesSchema) error {
	if r.rootClonesFile != "" {
		return fmt.Errorf(
			"%s: the clones configuration is already defined in %s",
			r.currentFile,
			r.rootClonesFile,
		)
	}

	c := Clones(s)

	if err := normalizePath(r.currentFile, &c.Dir); err != nil {
		return err
	}

	r.rootClonesFile = r.currentFile
	r.rootClones = c

	return nil
}

// populateRootClonesDefaults populates c with default values.
func (r *resolver) populateRootClonesDefaults(c *Clones) error {
	if c.Dir == "" {
		h, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf(
				"unable to determine default clones directory: %w",
				err,
			)
		}

		c.Dir = filepath.Join(h, "grit")
	}

	return nil
}

// normalizeSourceSpecificClonesBlock normalizes a clonesBlock within a source
// configuration.
func normalizeSourceSpecificClonesBlock(r *resolver, cfg unresolvedConfig, s *unresolvedSource) error {
	if s.Block.ClonesBlock == nil {
		s.Block.ClonesBlock = &clonesSchema{}
	}

	if s.Block.ClonesBlock.Dir == "" {
		s.Block.ClonesBlock.Dir = filepath.Join(
			r.rootClones.Dir,
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
