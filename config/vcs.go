package config

import (
	"fmt"
	"strings"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type unresolvedVCS struct {
	Block        vcsSchema
	DriverConfig vcsdriver.Config
	File         string
}

// mergeVCSDefaultsBlock merges a "vcs" block into the configuration.
func mergeVCSDefaultsBlock(
	reg *registry.Registry,
	cfg *unresolvedConfig,
	filename string,
	vcs vcsSchema,
) error {
	if vcs.DriverAlias == "" {
		return fmt.Errorf(
			"%s: this file contains a 'vcs' block with an empty driver alias",
			filename,
		)
	}

	for _, s := range cfg.VCSDefaults {
		if strings.EqualFold(s.Block.DriverAlias, vcs.DriverAlias) {
			return fmt.Errorf(
				"%s: defaults for the '%s' version control system are already defined in %s",
				filename,
				s.Block.DriverAlias,
				s.File,
			)
		}
	}

	r, ok := reg.VCSDriverByAlias(vcs.DriverAlias)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' version control system is not supported, the supported drivers are: '%s'",
			filename,
			vcs.DriverAlias,
			strings.Join(reg.VCSDriverAliases(), "', '"),
		)
	}

	driverBlock := r.NewConfigSchema()
	if diag := gohcl.DecodeBody(vcs.Body, nil, driverBlock); diag.HasErrors() {
		return diag
	}

	nc := &vcsNormalizationContext{filename}
	vcsConfig, err := driverBlock.NormalizeDefaults(nc)
	if err != nil {
		return fmt.Errorf(
			"%s: the default '%s' configuration is invalid: %w",
			filename,
			vcs.DriverAlias,
			err,
		)
	}

	if cfg.VCSDefaults == nil {
		cfg.VCSDefaults = map[string]unresolvedVCS{}
	}

	cfg.VCSDefaults[vcs.DriverAlias] = unresolvedVCS{
		Block:        vcs,
		DriverConfig: vcsConfig,
		File:         filename,
	}

	return nil
}

func mergeImplicitVCSDefaults(
	reg *registry.Registry,
	cfg *unresolvedConfig,
) error {
	if cfg.VCSDefaults == nil {
		cfg.VCSDefaults = map[string]unresolvedVCS{}
	}

	for alias, r := range reg.VCSDrivers() {
		if _, ok := cfg.VCSDefaults[alias]; ok { // TODO: case insensitive
			continue
		}

		nc := &vcsNormalizationContext{}
		dc, err := r.NewConfigSchema().NormalizeDefaults(nc)
		if err != nil {
			return err
		}

		cfg.VCSDefaults[alias] = unresolvedVCS{
			DriverConfig: dc,
		}
	}

	return nil
}

func normalizeSourceSpecificVCSBlocks(
	reg *registry.Registry,
	cfg unresolvedConfig,
	s *unresolvedSource,
) error {
	if s.VCSConfigs == nil {
		s.VCSConfigs = map[string]vcsdriver.Config{}
	}

	for _, b := range s.Block.VCSBlocks {
		if b.DriverAlias == "" {
			return fmt.Errorf(
				"%s: the '%s' source contains a 'vcs' block with an empty driver alias",
				s.File,
				s.Block.Name,
			)
		}

		for alias := range s.VCSConfigs {
			if strings.EqualFold(alias, b.DriverAlias) {
				return fmt.Errorf(
					"%s: the '%s' source contains multiple configurations for the '%s' version control system",
					s.File,
					s.Block.Name,
					b.DriverAlias,
				)
			}
		}

		r, ok := reg.VCSDriverByAlias(b.DriverAlias)
		if !ok {
			return fmt.Errorf(
				"%s: the '%s' source contains configuration for the '%s' version control system, which is not supported, the supported drivers are: '%s'",
				s.File,
				s.Block.Name,
				b.DriverAlias,
				strings.Join(reg.VCSDriverAliases(), "', '"),
			)
		}

		driverBlock := r.NewConfigSchema()
		if diag := gohcl.DecodeBody(b.Body, nil, driverBlock); diag.HasErrors() {
			return diag
		}

		var def vcsdriver.Config
		for alias, d := range cfg.VCSDefaults {
			if strings.EqualFold(alias, b.DriverAlias) {
				def = d.DriverConfig
				break
			}
		}

		nc := &vcsNormalizationContext{s.File}
		vcsConfig, err := driverBlock.NormalizeSourceSpecific(nc, def)
		if err != nil {
			return fmt.Errorf(
				"the '%s' configuration is invalid: %w",
				b.DriverAlias,
				err,
			)
		}

		s.VCSConfigs[b.DriverAlias] = vcsConfig
	}

	return nil
}

// vcsNormalizationContext is an implementation of
// vcsdriver.ConfigNormalizationContext.
type vcsNormalizationContext struct {
	filename string
}

func (c *vcsNormalizationContext) NormalizePath(p *string) error {
	// TODO: make resolution relative to the config directory, so that it's
	// always available.
	return normalizePath(c.filename, p)
}
