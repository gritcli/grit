package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// mergeGlobalVCS merges s into the configuration.
func (r *resolver) mergeGlobalVCS(s vcsSchema) error {
	if s.Driver == "" {
		return fmt.Errorf(
			"%s: global VCS configurations must provide a driver name",
			r.currentFile,
		)
	}

	if existingFile, ok := r.globalVCSFiles[s.Driver]; ok {
		return fmt.Errorf(
			"%s: global configuration for the '%s' version control system is already defined in %s",
			r.currentFile,
			s.Driver,
			existingFile,
		)
	}

	reg, ok := r.registry.VCSDriverByAlias(s.Driver)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' version control system is not supported, the supported VCS drivers are: '%s'",
			r.currentFile,
			s.Driver,
			strings.Join(r.registry.VCSDriverAliases(), "', '"),
		)
	}

	bodySchema := reg.NewConfigSchema()
	if diag := gohcl.DecodeBody(s.Body, nil, bodySchema); diag.HasErrors() {
		return diag
	}

	nc := &normalizeContext{r.currentFile}
	cfg, err := bodySchema.NormalizeGlobals(nc)
	if err != nil {
		return fmt.Errorf(
			"%s: the global configuration for the '%s' version control system is invalid: %w",
			r.currentFile,
			s.Driver,
			err,
		)
	}

	if r.globalVCSs == nil {
		r.globalVCSFiles = map[string]string{}
		r.globalVCSs = map[string]vcsdriver.Config{}
	}

	r.globalVCSFiles[s.Driver] = r.currentFile
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

		nc := &normalizeContext{r.configDir}
		cfg, err := reg.NewConfigSchema().NormalizeGlobals(nc)
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

func normalizeSourceSpecificVCSBlocks(
	r *resolver,
	cfg unresolvedConfig,
	s *unresolvedSource,
) error {
	if s.VCSConfigs == nil {
		s.VCSConfigs = map[string]vcsdriver.Config{}
	}

	for _, b := range s.Block.VCSBlocks {
		if b.Driver == "" {
			return fmt.Errorf(
				"%s: the '%s' source contains a 'vcs' block with an empty driver alias",
				s.File,
				s.Block.Name,
			)
		}

		for alias := range s.VCSConfigs {
			if alias == b.Driver {
				return fmt.Errorf(
					"%s: the '%s' source contains multiple configurations for the '%s' version control system",
					s.File,
					s.Block.Name,
					b.Driver,
				)
			}
		}

		reg, ok := r.registry.VCSDriverByAlias(b.Driver)
		if !ok {
			return fmt.Errorf(
				"%s: the '%s' source contains configuration for the '%s' version control system, which is not supported, the supported drivers are: '%s'",
				s.File,
				s.Block.Name,
				b.Driver,
				strings.Join(r.registry.VCSDriverAliases(), "', '"),
			)
		}

		driverBlock := reg.NewConfigSchema()
		if diag := gohcl.DecodeBody(b.Body, nil, driverBlock); diag.HasErrors() {
			return diag
		}

		nc := &normalizeContext{filepath.Dir(s.File)}
		vcsConfig, err := driverBlock.NormalizeSourceSpecific(nc, r.globalVCSs[b.Driver])
		if err != nil {
			return fmt.Errorf(
				"the '%s' configuration is invalid: %w",
				b.Driver,
				err,
			)
		}

		s.VCSConfigs[b.Driver] = vcsConfig
	}

	return nil
}
