package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
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
	Driver      string       `hcl:",label"`
	Enabled     *bool        `hcl:"enabled"`
	ClonesBlock *clonesBlock `hcl:"clones,block"`
	DriverBlock hcl.Body     `hcl:",remain"` // parsed into a sourceDriverBlock, as per sourceDriverBlockFactory
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

	d, ok := reg.SourceDriverByAlias(b.Driver)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' source uses '%s' which is not supported, the supported drivers are: '%s'",
			filename,
			b.Name,
			b.Driver,
			strings.Join(reg.SourceDriverAliases(), "', '"),
		)
	}

	driverBlock := d.NewConfigSchema()
	if diag := gohcl.DecodeBody(b.DriverBlock, nil, driverBlock); diag.HasErrors() {
		return diag
	}

	if cfg.Sources == nil {
		cfg.Sources = map[string]unresolvedSource{}
	}

	cfg.Sources[b.Name] = unresolvedSource{
		File:        filename,
		Block:       b,
		DriverBlock: driverBlock,
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
func normalizeSourceBlock(cfg unresolvedConfig, s *unresolvedSource) error {
	if s.Block.Enabled == nil {
		enabled := true
		s.Block.Enabled = &enabled
	}

	return normalizeSourceSpecificClonesBlock(cfg, s)
}

// assembleSourceBlock converts b into its configuration representation.
func assembleSourceBlock(cfg unresolvedConfig, s unresolvedSource) (Source, error) {
	nc := &normalizationContext{cfg, s}

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

// normalizationContext is an implementation of
// sourcedriver.ConfigNormalizationContext.
type normalizationContext struct {
	cfg unresolvedConfig
	s   unresolvedSource
}

func (c *normalizationContext) NormalizePath(p *string) error {
	return normalizePath(c.s.File, p)
}

func (c *normalizationContext) ResolveVCSConfig(in, out interface{}) error {
	switch out := out.(type) {
	case *Git:
		b := in.(*gitBlock)
		if err := normalizeSourceSpecificGitBlock(c.cfg, c.s, &b); err != nil {
			return err
		}

		*out = assembleGitBlock(*b)

	default:
		return fmt.Errorf("unsupported VCS config type (%T)", out)
	}

	return nil
}
