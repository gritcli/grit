package config

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// mergeGlobalClones merges s into the configuration.
func (l *loader) mergeGlobalClones(file string, s clonesSchema) error {
	if l.globalClonesFile != "" {
		return fmt.Errorf(
			"%s: the global clones configuration is already defined in %s",
			file,
			l.globalClonesFile,
		)
	}

	cfg := Clones(s)

	if err := l.normalizePath(&cfg.Dir); err != nil {
		return err
	}

	l.globalClonesFile = file
	l.globalClones = cfg

	return nil
}

// populateGlobalClonesDefaults populates l.globalClones with default values.
func (l *loader) populateGlobalClonesDefaults() error {
	if l.globalClones.Dir == "" {
		h, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf(
				"unable to determine default clones directory: %w",
				err,
			)
		}

		l.globalClones.Dir = filepath.Join(h, "grit")
	}

	return nil
}

// finalizeSouceSpecific returns the clones configuration to use for a specific
// source.
func (l *loader) finalizeSourceSpecificClones(
	i intermediateSource,
	s *clonesSchema,
) (Clones, error) {
	cfg := Clones{}

	if s != nil {
		cfg.Dir = s.Dir

		if err := l.normalizePath(&cfg.Dir); err != nil {
			return Clones{}, err
		}
	}

	if cfg.Dir == "" {
		cfg.Dir = filepath.Join(l.globalClones.Dir, i.Schema.Name)
	}

	return cfg, nil
}
