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

	ctx := &vcsContext{l}
	cfg, err := reg.ConfigLoader.UnmarshalAndMerge(
		ctx,
		l.defaultVCSs[s.Driver],
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

// populateDefaultVCSs populates l.defaultVCSs with default configurations for
// each of the supported VCS drivers.
func (l *loader) populateDefaultVCSs() error {
	l.defaultVCSs = map[string]vcsdriver.Config{}

	for alias, reg := range l.Registry.VCSDrivers() {
		ctx := &vcsContext{l}
		cfg, err := reg.ConfigLoader.Defaults(ctx)
		if err != nil {
			return fmt.Errorf(
				"unable to produce default global configuration for the '%s' version control system: %w",
				alias,
				err,
			)
		}

		l.defaultVCSs[alias] = cfg
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

		cfg, ok := l.globalVCSs[driver]
		if !ok {
			cfg = l.defaultVCSs[driver]
		}

		ctx := &vcsContext{l}
		cfg, err := reg.ConfigLoader.UnmarshalAndMerge(ctx, cfg, body)
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

// vcsContext is an implementation of the vcsdriver.ConfigContext interface.
type vcsContext struct {
	loader *loader
}

func (c *vcsContext) EvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{}
}

func (c *vcsContext) NormalizePath(p *string) error {
	return c.loader.normalizePath(p)
}
