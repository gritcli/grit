package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// Source represents a repository source defined in the configuration.
type Source struct {
	// Name is a short identifier for the source. Each source in the
	// configuration has a unique name. Names are case-insensitive.
	Name string

	// Enabled is true if the source is enabled. Disabled sources are not used
	// when searching for repositories to be cloned.
	Enabled bool

	// Clones is the configuration that controls how Grit stores local
	// repository clones for this source.
	Clones Clones

	// Driver contains driver-specific configuration for this source.
	Driver sourcedriver.Config
}

// sourceBlock is the HCL schema for a "source" block.
type sourceBlock struct {
	Name        string       `hcl:",label"`
	DriverAlias string       `hcl:",label"`
	Enabled     *bool        `hcl:"enabled"`
	ClonesBlock *clonesBlock `hcl:"clones,block"`
	VCSBlocks   []vcsBlock   `hcl:"vcs,block"`
	DriverBlock hcl.Body     `hcl:",remain"`
}

// unresolvedSource contains information about a "source" block within an
// as-yet-unresolved configuration.
type unresolvedSource struct {
	Block       sourceBlock
	DriverBlock sourcedriver.ConfigSchema
	VCSConfigs  map[string]vcsdriver.Config
	File        string
}

// sourceNameRegexp is a regular expression used to validate source names.
var sourceNameRegexp = regexp.MustCompile(`(?i)^[a-z_]+$`)

// mergeSourceBlock merges b into cfg.
func mergeSourceBlock(
	reg *registry.Registry,
	cfg *unresolvedConfig,
	filename string,
	b sourceBlock,
) error {
	if b.Name == "" {
		return fmt.Errorf(
			"%s: this file contains a 'source' block with an empty name",
			filename,
		)
	}

	if !sourceNameRegexp.MatchString(b.Name) {
		return fmt.Errorf(
			"%s: the '%s' source has an invalid name, source names must contain only alpha-numeric characters and underscores",
			filename,
			b.Name,
		)
	}

	for _, s := range cfg.Sources {
		if strings.EqualFold(s.Block.Name, b.Name) {
			return fmt.Errorf(
				"%s: a source named '%s' is already defined in %s",
				filename,
				s.Block.Name,
				s.File,
			)
		}
	}

	r, ok := reg.SourceDriverByAlias(b.DriverAlias)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' source uses '%s' which is not supported, the supported drivers are: '%s'",
			filename,
			b.Name,
			b.DriverAlias,
			strings.Join(reg.SourceDriverAliases(), "', '"),
		)
	}

	driverBlock := r.NewConfigSchema()
	if diag := gohcl.DecodeBody(b.DriverBlock, nil, driverBlock); diag.HasErrors() {
		return diag
	}

	if cfg.Sources == nil {
		cfg.Sources = map[string]unresolvedSource{}
	}

	cfg.Sources[b.Name] = unresolvedSource{
		Block:       b,
		DriverBlock: driverBlock,
		File:        filename,
	}

	return nil
}

// mergeDefaultSources merges any default sources into cfg, if it does not
// already contain a source with the same name.
func mergeDefaultSources(
	reg *registry.Registry,
	cfg *unresolvedConfig,
) {
	if cfg.Sources == nil {
		cfg.Sources = map[string]unresolvedSource{}
	}

	for _, alias := range reg.SourceDriverAliases() {
		reg, _ := reg.SourceDriverByAlias(alias)

		for n, new := range reg.DefaultSources {
			if _, ok := cfg.Sources[n]; ok {
				continue
			}

			cfg.Sources[n] = unresolvedSource{
				Block: sourceBlock{
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
	reg *registry.Registry,
	cfg unresolvedConfig,
	s *unresolvedSource,
) error {
	if s.Block.Enabled == nil {
		enabled := true
		s.Block.Enabled = &enabled
	}

	if err := normalizeSourceSpecificClonesBlock(cfg, s); err != nil {
		return err
	}

	return normalizeSourceSpecificVCSBlocks(reg, cfg, s)
}

// assembleSourceBlock converts b into its configuration representation.
func assembleSourceBlock(cfg unresolvedConfig, s unresolvedSource) (Source, error) {
	nc := &sourceNormalizationContext{cfg, s}

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
		Clones:  assembleClonesBlock(*s.Block.ClonesBlock),
		Driver:  driverConfig,
	}, nil
}

// sourceNormalizationContext is an implementation of
// sourcedriver.ConfigNormalizationContext.
type sourceNormalizationContext struct {
	cfg unresolvedConfig
	s   unresolvedSource
}

func (c *sourceNormalizationContext) NormalizePath(p *string) error {
	// TODO: make resolution relative to the config directory, so that it's
	// always available.
	return normalizePath(c.s.File, p)
}

func (c *sourceNormalizationContext) ResolveVCSConfig(cfg interface{}) error {
	elem := reflect.ValueOf(cfg).Elem()

	for _, d := range c.s.VCSConfigs {
		v := reflect.ValueOf(d)

		if v.Type().AssignableTo(
			elem.Type(),
		) {
			elem.Set(v)
			return nil
		}
	}

	for _, d := range c.cfg.VCSDefaults {
		v := reflect.ValueOf(d.DriverConfig)

		if v.Type().AssignableTo(
			elem.Type(),
		) {
			elem.Set(v)
			return nil
		}
	}

	return fmt.Errorf("unsupported VCS config type (%s)", elem.Type())
}
