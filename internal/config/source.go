package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

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

	// DriverConfig contains driver-specific configuration for this source.
	DriverConfig SourceDriverConfig
}

// AcceptVisitor calls the appropriate method on v.
func (s Source) AcceptVisitor(v SourceVisitor) {
	s.DriverConfig.acceptVisitor(s, v)
}

// SourceDriverConfig is an interface for driver-specific configuration options
// for a repository source.
type SourceDriverConfig interface {
	// acceptVisitor calls the appropriate method on v.
	acceptVisitor(s Source, v SourceVisitor)
}

// SourceVisitor dispatches Source values to driver-specific logic.
type SourceVisitor interface {
	VisitGitHubSource(s Source, cfg GitHub)
}

// sourceBlock is the HCL schema for a "source" block.
type sourceBlock struct {
	Name           string       `hcl:",label"`
	Driver         string       `hcl:",label"`
	Enabled        *bool        `hcl:"enabled"`
	ClonesBlock    *clonesBlock `hcl:"clones,block"`
	DriverSpecific hcl.Body     `hcl:",remain"` // parsed into a sourceDriverBlock, as per sourceDriverBlockFactory
}

// sourceDriverBlock is an interface for the HCL schema for driver-specific
// parts of a "source" block's body.
type sourceDriverBlock interface {
	// Normalize normalizes the block in-place.
	Normalize(cfg unresolvedConfig, s unresolvedSource) error

	// Assemble converts the block into its configuration representation.
	Assemble() SourceDriverConfig
}

// sourceNameRegexp is a regular expression used to validate source names.
var sourceNameRegexp = regexp.MustCompile(`(?i)^[a-z_]+$`)

// mergeSourceBlock merges b into cfg.
func mergeSourceBlock(cfg *unresolvedConfig, filename string, b sourceBlock) error {
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
				b.Name,
				s.File,
			)
		}
	}

	newBody, ok := sourceDriverBlockFactory[b.Driver]
	if !ok {
		var drivers []string
		for n := range sourceDriverBlockFactory {
			drivers = append(drivers, n)
		}
		sort.Strings(drivers)

		return fmt.Errorf(
			"%s: the '%s' source uses '%s' which is not supported, the supported drivers are: '%s'",
			filename,
			b.Name,
			b.Driver,
			strings.Join(drivers, "', '"),
		)
	}

	body := newBody()
	if diag := gohcl.DecodeBody(b.DriverSpecific, nil, body); diag.HasErrors() {
		return diag
	}

	if cfg.Sources == nil {
		cfg.Sources = map[string]unresolvedSource{}
	}

	cfg.Sources[b.Name] = unresolvedSource{
		File:  filename,
		Block: b,
		Body:  body,
	}

	return nil
}

// mergeDefaultSources merges any default sources into cfg, if it does not
// already contain a source with the same name.
func mergeDefaultSources(cfg *unresolvedConfig) {
	if cfg.Sources == nil {
		cfg.Sources = map[string]unresolvedSource{}
	}

	for n, newBody := range defaultSourceFactoryByName {
		if cfg.Sources == nil {
			cfg.Sources = map[string]unresolvedSource{}
		} else if _, ok := cfg.Sources[n]; ok {
			continue
		}

		cfg.Sources[n] = unresolvedSource{
			Block: sourceBlock{
				Name: n,
			},
			Body: newBody(),
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

	if err := s.Body.Normalize(cfg, *s); err != nil {
		return fmt.Errorf(
			"%s: the '%s' repository source is invalid: %w",
			s.File,
			s.Block.Name,
			err,
		)
	}

	return normalizeSourceSpecificClonesBlock(cfg, s)
}

// assembleSourceBlock converts b into its configuration representation.
func assembleSourceBlock(b sourceBlock, body sourceDriverBlock) Source {
	return Source{
		Name:         b.Name,
		Enabled:      *b.Enabled,
		Clones:       assembleClonesBlock(*b.ClonesBlock),
		DriverConfig: body.Assemble(),
	}
}

var (
	// sourceDriverBlockFactory is a map of a source driver name to a function
	// that returns a new, empty sourceDriverBlock for that driver.
	sourceDriverBlockFactory = map[string]func() sourceDriverBlock{}

	// defaultSourceFactoryByName is a map of a source name to a function that
	// returns a new default source. These defaults are merged into any
	// configuration that does not already contain a repository source with the
	// same name.
	defaultSourceFactoryByName = map[string]func() sourceDriverBlock{}
)

// registerSourceDriver registers a source driver, allowing its configuration to
// be parsed.
//
// name is the name of the driver, which is given as the second "label" (HCL
// terminology) on the "source" blocks within the configuration file.
func registerSourceDriver(
	name string,
	newBlock func() sourceDriverBlock,
) {
	if _, ok := sourceDriverBlockFactory[name]; ok {
		panic("source driver name already registered")
	}

	sourceDriverBlockFactory[name] = newBlock
}

// registerDefaultSource registers a default source that is merged into every
// configuration unless overridden by the user.
func registerDefaultSource(
	name string,
	newBody func() sourceDriverBlock,
) {
	if _, ok := defaultSourceFactoryByName[name]; ok {
		panic("default source name already registered")
	}

	defaultSourceFactoryByName[name] = newBody
}
