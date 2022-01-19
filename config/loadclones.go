package config

import (
	"fmt"
	"path/filepath"
)

// mergeGlobalClones merges s into the configuration.
func (l *loader) mergeGlobalClones(file string, s clonesSchema) error {
	if l.globalClonesFile != "" {
		return fmt.Errorf(
			"the global clones configuration is already defined in %s",
			l.globalClonesFile,
		)
	}

	cfg := Clones(s)

	if err := l.normalizePath(&cfg.Dir); err != nil {
		return fmt.Errorf(
			"unable to resolve global clones directory: %w (%s)",
			err,
			cfg.Dir,
		)
	}

	l.globalClonesFile = file
	l.globalClones = cfg

	return nil
}

// populateGlobalClonesDefaults populates l.globalClones with default values.
func (l *loader) populateGlobalClonesDefaults() error {
	if l.globalClones.Dir == "" {
		l.globalClones.Dir = DefaultClonesDirectory

		if err := l.normalizePath(&l.globalClones.Dir); err != nil {
			return fmt.Errorf(
				"unable to resolve default global clones directory: %w (%s)",
				err,
				l.globalClones.Dir,
			)
		}
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
			return Clones{}, fmt.Errorf(
				"unable to resolve clones directory for the '%s' source: %w (%s)",
				i.Schema.Name,
				err,
				cfg.Dir,
			)
		}
	}

	if cfg.Dir == "" {
		cfg.Dir = filepath.Join(l.globalClones.Dir, i.Schema.Name)
	}

	return cfg, nil
}
