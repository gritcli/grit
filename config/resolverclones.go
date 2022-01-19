package config

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// mergeGlobalClones merges s into the configuration.
func (r *resolver) mergeGlobalClones(file string, s clonesSchema) error {
	if r.globalClonesFile != "" {
		return fmt.Errorf(
			"%s: the global clones configuration is already defined in %s",
			file,
			r.globalClonesFile,
		)
	}

	cfg := Clones(s)

	if err := r.normalizePath(&cfg.Dir); err != nil {
		return err
	}

	r.globalClonesFile = file
	r.globalClones = cfg

	return nil
}

// populateGlobalClonesDefaults populates r.globalClones with default values.
// TODO: can this be moved into the mergeGlobalClones() function
func (r *resolver) populateGlobalClonesDefaults() error {
	if r.globalClones.Dir == "" {
		h, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf(
				"unable to determine default clones directory: %w",
				err,
			)
		}

		r.globalClones.Dir = filepath.Join(h, "grit")
	}

	return nil
}

// finalizeSouceSpecific returns the clones configuration to use for a specific
// source.
func (r *resolver) finalizeSourceSpecificClones(
	i intermediateSource,
	s *clonesSchema,
) (Clones, error) {
	cfg := Clones{}

	if s != nil {
		cfg.Dir = s.Dir

		if err := r.normalizePath(&cfg.Dir); err != nil {
			return Clones{}, err
		}
	}

	if cfg.Dir == "" {
		cfg.Dir = filepath.Join(r.globalClones.Dir, i.Schema.Name)
	}

	return cfg, nil
}
