package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2/gohcl"
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

	bodySchema := reg.NewConfigSchema()
	if diag := gohcl.DecodeBody(s.DriverBody, nil, bodySchema); diag.HasErrors() {
		return diag
	}

	cfg, err := bodySchema.NormalizeGlobals(&vcsNormalizeContext{l})
	if err != nil {
		return fmt.Errorf(
			"the global configuration for the '%s' version control system can not be loaded: %w",
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

		cfg, err := reg.NewConfigSchema().NormalizeGlobals(&vcsNormalizeContext{l})
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

	reg, ok := l.Registry.VCSDriverByAlias(s.Driver)
	if !ok {
		return fmt.Errorf(
			"the '%s' source contains configuration for an unrecognized version control system ('%s'), the supported VCS drivers are '%s'",
			i.Schema.Name,
			s.Driver,
			strings.Join(l.Registry.VCSDriverAliases(), "', '"),
		)
	}

	bodySchema := reg.NewConfigSchema()
	if diag := gohcl.DecodeBody(s.DriverBody, nil, bodySchema); diag.HasErrors() {
		return diag
	}

	i.VCSs[s.Driver] = bodySchema

	return nil
}

// finalizeSourceSpecificVCS returns the VCS configurations for a specific
// source.
func (l *loader) finalizeSourceSpecificVCSs(
	i intermediateSource,
) (map[string]vcsdriver.Config, error) {
	configs := map[string]vcsdriver.Config{}

	for n, s := range i.VCSs {
		cfg, err := s.NormalizeSourceSpecific(
			&vcsNormalizeContext{l},
			l.globalVCSs[n],
		)
		if err != nil {
			return nil, fmt.Errorf(
				"the '%s' source's configuration for the '%s' version control system could not be loaded: %w",
				i.Schema.Name,
				n,
				err,
			)
		}

		configs[n] = cfg
	}

	return configs, nil
}

// vcsNormalizeContext is an implementation of the
// vcsdriver.ConfigNormalizeContext interface.
type vcsNormalizeContext struct {
	loader *loader
}

func (nc *vcsNormalizeContext) NormalizePath(p *string) error {
	return nc.loader.normalizePath(p)
}
