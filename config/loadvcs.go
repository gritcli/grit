package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2"
)

// mergeGlobalVCS merges s into the configuration.
func (l *loader) mergeGlobalVCS(file string, s vcsSchema) error {
	if s.Driver == "" {
		return errors.New("global VCS configuration with empty driver name")
	}

	if existingFile, ok := l.globalVCSFiles[s.Driver]; ok {
		return fmt.Errorf(
			"global configuration for the '%s' version control system is already defined in %s",
			s.Driver,
			existingFile,
		)
	}

	reg, ok := l.Registry.VCSDriverByAlias(s.Driver)
	if !ok {
		return fmt.Errorf(
			"the '%s' version control system is not unrecognized, the supported VCS drivers are: '%s'",
			s.Driver,
			strings.Join(l.Registry.VCSDriverAliases(), "', '"),
		)
	}

	nc := &vcsNormalizeContext{l}

	cfg, err := reg.ConfigNormalizer.Defaults(nc)
	if err != nil {
		if isHCLError(err) {
			return err
		}

		return fmt.Errorf(
			"the default configuration for the '%s' version control system cannot be loaded: %w",
			s.Driver,
			err,
		)
	}

	cfg, err = reg.ConfigNormalizer.Merge(
		nc,
		cfg,
		s.DriverBody,
	)
	if err != nil {
		if isHCLError(err) {
			return err
		}

		return fmt.Errorf(
			"the global configuration for the '%s' version control system cannot be loaded: %w",
			s.Driver,
			err,
		)
	}

	if l.globalVCSs == nil {
		l.globalVCSFiles = map[string]string{}
		l.globalVCSs = map[string]vcsdriver.Config{}
	}

	l.globalVCSFiles[s.Driver] = file
	l.globalVCSs[s.Driver] = cfg

	return nil
}

// populateGlobalClonesDefaults populates l.globalVCSs with default
// configurations for each of the supported VCS drivers.
func (l *loader) populateImplicitGlobalVCSs() error {
	if l.globalVCSs == nil {
		l.globalVCSs = map[string]vcsdriver.Config{}
	}

	for alias, reg := range l.Registry.VCSDrivers() {
		if _, ok := l.globalVCSs[alias]; ok {
			continue
		}

		nc := &vcsNormalizeContext{l}
		cfg, err := reg.ConfigNormalizer.Defaults(nc)
		if err != nil {
			return fmt.Errorf(
				"unable to produce default global configuration for the '%s' version control system: %w",
				alias,
				err,
			)
		}

		l.globalVCSs[alias] = cfg
	}

	return nil
}

// mergeSourceSpecificVCS merges s into i.VCSs.
func (l *loader) mergeSourceSpecificVCS(i *intermediateSource, s vcsSchema) error {
	if s.Driver == "" {
		return fmt.Errorf(
			"the '%s' source contains a VCS configuration with an empty driver name",
			i.Schema.Name,
		)
	}

	if _, ok := i.VCSs[s.Driver]; ok {
		return fmt.Errorf(
			"the '%s' source contains multiple configurations for the '%s' version control system",
			i.Schema.Name,
			s.Driver,
		)
	}

	i.VCSs[s.Driver] = s.DriverBody

	return nil
}

// finalizeSourceSpecificVCS returns the VCS configurations for a specific
// source.
func (l *loader) finalizeSourceSpecificVCSs(
	i intermediateSource,
) (map[string]vcsdriver.Config, error) {
	configs := map[string]vcsdriver.Config{}

	for driver, body := range i.VCSs {
		reg, ok := l.Registry.VCSDriverByAlias(driver)
		if !ok {
			return nil, fmt.Errorf(
				"the '%s' source contains configuration for an unrecognized version control system ('%s'), the supported VCS drivers are '%s'",
				i.Schema.Name,
				driver,
				strings.Join(l.Registry.VCSDriverAliases(), "', '"),
			)
		}

		nc := &vcsNormalizeContext{l}
		cfg, err := reg.ConfigNormalizer.Merge(nc, l.globalVCSs[driver], body)
		if err != nil {
			if isHCLError(err) {
				return nil, err
			}

			return nil, fmt.Errorf(
				"the '%s' source's configuration for the '%s' version control system cannot be loaded: %w",
				i.Schema.Name,
				driver,
				err,
			)
		}

		configs[driver] = cfg
	}

	return configs, nil
}

// vcsNormalizeContext is an implementation of the
// vcsdriver.ConfigNormalizeContext interface.
type vcsNormalizeContext struct {
	loader *loader
}

func (nc *vcsNormalizeContext) EvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{}
}

func (nc *vcsNormalizeContext) NormalizePath(p *string) error {
	return nc.loader.normalizePath(p)
}
