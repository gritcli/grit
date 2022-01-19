package config

import (
	"fmt"
	"strings"

	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// mergeGlobalVCS merges s into the configuration.
func (r *resolver) mergeGlobalVCS(file string, s vcsSchema) error {
	if s.Driver == "" {
		return fmt.Errorf(
			"%s: global VCS configuration with empty driver name",
			file,
		)
	}

	if existingFile, ok := r.globalVCSFiles[s.Driver]; ok {
		return fmt.Errorf(
			"%s: global configuration for the '%s' version control system is already defined in %s",
			file,
			s.Driver,
			existingFile,
		)
	}

	reg, ok := r.registry.VCSDriverByAlias(s.Driver)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' version control system is not unrecognized, the supported VCS drivers are: '%s'",
			file,
			s.Driver,
			strings.Join(r.registry.VCSDriverAliases(), "', '"),
		)
	}

	bodySchema := reg.NewConfigSchema()
	if diag := gohcl.DecodeBody(s.DriverBody, nil, bodySchema); diag.HasErrors() {
		return diag
	}

	cfg, err := bodySchema.NormalizeGlobals(&vcsNormalizeContext{r})
	if err != nil {
		return fmt.Errorf(
			"%s: the global configuration for the '%s' version control system is invalid: %w",
			file,
			s.Driver,
			err,
		)
	}

	if r.globalVCSs == nil {
		r.globalVCSFiles = map[string]string{}
		r.globalVCSs = map[string]vcsdriver.Config{}
	}

	r.globalVCSFiles[s.Driver] = file
	r.globalVCSs[s.Driver] = cfg

	return nil
}

// populateGlobalClonesDefaults populates r.globalVCSs with default
// configurations for each of the supported VCS drivers.
func (r *resolver) populateImplicitGlobalVCSs() error {
	if r.globalVCSs == nil {
		r.globalVCSs = map[string]vcsdriver.Config{}
	}

	for alias, reg := range r.registry.VCSDrivers() {
		if _, ok := r.globalVCSs[alias]; ok {
			continue
		}

		cfg, err := reg.NewConfigSchema().NormalizeGlobals(&vcsNormalizeContext{r})
		if err != nil {
			return fmt.Errorf(
				"unable to produce default global configuration for the '%s' version control system: %w",
				alias,
				err,
			)
		}

		r.globalVCSs[alias] = cfg
	}

	return nil
}

// mergeSourceSpecificVCS merges s into i.VCSs.
func (r *resolver) mergeSourceSpecificVCS(i *intermediateSource, s vcsSchema) error {
	if s.Driver == "" {
		return fmt.Errorf(
			"%s: the '%s' source contains a VCS configuration with an empty driver name",
			i.File,
			i.Schema.Name,
		)
	}

	if _, ok := i.VCSs[s.Driver]; ok {
		return fmt.Errorf(
			"%s: the '%s' source contains multiple configurations for the '%s' version control system",
			i.File,
			i.Schema.Name,
			s.Driver,
		)
	}

	reg, ok := r.registry.VCSDriverByAlias(s.Driver)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' source contains configuration for an unrecognized version control system ('%s'), the supported VCS drivers are '%s'",
			i.File,
			i.Schema.Name,
			s.Driver,
			strings.Join(r.registry.VCSDriverAliases(), "', '"),
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
func (r *resolver) finalizeSourceSpecificVCSs(
	i intermediateSource,
) (map[string]vcsdriver.Config, error) {
	configs := map[string]vcsdriver.Config{}

	for n, s := range i.VCSs {
		cfg, err := s.NormalizeSourceSpecific(
			&vcsNormalizeContext{r},
			r.globalVCSs[n],
		)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: the '%s' source's configuration for the '%s' version control system is invalid: %w",
				i.File, // TODO: this will be empty for 'implicit' sources
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
	resolver *resolver
}

func (nc *vcsNormalizeContext) NormalizePath(p *string) error {
	return nc.resolver.normalizePath(p)
}
