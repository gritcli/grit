package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// unresolvedSource contains information about a "source" block within an
// as-yet-unresolved configuration.
type unresolvedSource struct {
	Block       sourceSchema
	DriverBlock sourcedriver.ConfigSchema
	VCSConfigs  map[string]vcsdriver.Config
	File        string
}

// sourceNameRegexp is a regular expression used to validate source names.
var sourceNameRegexp = regexp.MustCompile(`(?i)^[a-z_]+$`)

// mergeSourceBlock merges a "source" block into the configuration.
func mergeSourceBlock(
	reg *registry.Registry,
	cfg *unresolvedConfig,
	filename string,
	src sourceSchema,
) error {
	if src.Name == "" {
		return fmt.Errorf(
			"%s: this file contains a 'source' block with an empty name",
			filename,
		)
	}

	if !sourceNameRegexp.MatchString(src.Name) {
		return fmt.Errorf(
			"%s: the '%s' source has an invalid name, source names must contain only alpha-numeric characters and underscores",
			filename,
			src.Name,
		)
	}

	for _, s := range cfg.Sources {
		if strings.EqualFold(s.Block.Name, src.Name) {
			return fmt.Errorf(
				"%s: a source named '%s' is already defined in %s",
				filename,
				s.Block.Name,
				s.File,
			)
		}
	}

	r, ok := reg.SourceDriverByAlias(src.DriverAlias)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' source uses '%s' which is not supported, the supported drivers are: '%s'",
			filename,
			src.Name,
			src.DriverAlias,
			strings.Join(reg.SourceDriverAliases(), "', '"),
		)
	}

	driverBlock := r.NewConfigSchema()
	if diag := gohcl.DecodeBody(src.DriverBlock, nil, driverBlock); diag.HasErrors() {
		return diag
	}

	if cfg.Sources == nil {
		cfg.Sources = map[string]unresolvedSource{}
	}

	cfg.Sources[src.Name] = unresolvedSource{
		Block:       src,
		DriverBlock: driverBlock,
		File:        filename,
	}

	return nil
}

// mergeImplicitSources merges any implicit sources into cfg, if it does not
// already contain a source with the same name.
func mergeImplicitSources(
	reg *registry.Registry,
	cfg *unresolvedConfig,
) {
	if cfg.Sources == nil {
		cfg.Sources = map[string]unresolvedSource{}
	}

	for _, alias := range reg.SourceDriverAliases() {
		reg, _ := reg.SourceDriverByAlias(alias)

		for n, new := range reg.ImplicitSources {
			if _, ok := cfg.Sources[n]; ok {
				continue
			}

			cfg.Sources[n] = unresolvedSource{
				Block: sourceSchema{
					Name: n,
				},
				DriverBlock: new(),
			}
		}
	}
}

// normalizeSourceBlock normalizes cfg.Sources and populates them with
// default values.
func normalizeSourceBlock(
	r *resolver,
	reg *registry.Registry,
	cfg unresolvedConfig,
	s *unresolvedSource,
) error {
	if s.Block.Enabled == nil {
		enabled := true
		s.Block.Enabled = &enabled
	}

	if err := normalizeSourceSpecificClonesBlock(r, cfg, s); err != nil {
		return err
	}

	return normalizeSourceSpecificVCSBlocks(r, cfg, s)
}

// assembleSourceBlock converts b into its configuration representation.
func assembleSourceBlock(r *resolver, cfg unresolvedConfig, s unresolvedSource) (Source, error) {
	nc := &sourceNormalizationContext{r.globalVCSs, cfg, s}

	// TODO: fail if Normalize() does not call nc.ResolveVCSConfig() for every
	// type of VCS config specified on this source.
	driverConfig, err := s.DriverBlock.Normalize(nc)
	if err != nil {
		return Source{}, fmt.Errorf(
			"%s: the '%s' repository source is invalid: %w",
			s.File,
			s.Block.Name,
			err,
		)
	}

	return Source{
		Name:    s.Block.Name,
		Enabled: *s.Block.Enabled,
		Clones:  Clones(*s.Block.ClonesBlock),
		Driver:  driverConfig,
	}, nil
}

// sourceNormalizationContext is an implementation of
// sourcedriver.ConfigNormalizationContext.
type sourceNormalizationContext struct {
	globalVCSs map[string]vcsdriver.Config
	cfg        unresolvedConfig
	s          unresolvedSource
}

func (c *sourceNormalizationContext) NormalizePath(p *string) error {
	// TODO: make resolution relative to the config directory, so that it's
	// always available.
	return normalizePath(c.s.File, p)
}

func (c *sourceNormalizationContext) ResolveVCSConfig(cfg interface{}) error {
	// TODO: validate that cfg is a pointer-to-impl to provide a better error
	// message.
	elem := reflect.ValueOf(cfg).Elem()

	if c.resolveVCSConfig(c.s.VCSConfigs, elem) {
		return nil
	}

	if c.resolveVCSConfig(c.globalVCSs, elem) {
		return nil
	}

	return fmt.Errorf(
		"none of the supported VCS drivers provided a config of type %s",
		elem.Type(),
	)
}

// resolveVCSConfig assigns a value from configs to elem, if possible.
func (c *sourceNormalizationContext) resolveVCSConfig(
	configs map[string]vcsdriver.Config,
	elem reflect.Value,
) bool {
	t := elem.Type()

	for _, cfg := range configs {
		v := reflect.ValueOf(cfg)

		if v.Type().AssignableTo(t) {
			elem.Set(v)
			return true
		}
	}

	return false
}
