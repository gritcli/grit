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
	// configuration has a unique name.
	Name string

	// Enabled is true if the source is enabled. Disabled sources are not used
	// when searching for repositories to be cloned.
	Enabled bool

	// Config contains implementation-specific configuration for this source.
	Config SourceConfig
}

// AcceptVisitor calls the appropriate method on v.
func (s Source) AcceptVisitor(v SourceVisitor) {
	s.Config.acceptVisitor(s, v)
}

// SourceConfig is an interface for implementation-specific configuration
// options for a repository source.
type SourceConfig interface {
	// acceptVisitor calls the appropriate method on v.
	acceptVisitor(s Source, v SourceVisitor)
}

// SourceVisitor dispatches Source values to implementation-specific logic.
type SourceVisitor interface {
	VisitGitHubSource(s Source, cfg GitHubConfig)
}

// sourceBlock is the HCL schema for a "source" block.
type sourceBlock struct {
	Name    string   `hcl:",label"`
	Impl    string   `hcl:",label"`
	Enabled *bool    `hcl:"enabled"`
	Body    hcl.Body `hcl:",remain"` // see sourceBlockBody
}

// sourceBlockBody is an interface for implementation-specific HCL schema within
// the body of a source block.
type sourceBlockBody interface {
	// Normalize normalizes the body in-place.
	Normalize(cfg unresolvedConfig) error

	// Assemble converts the body into its configuration representation.
	Assemble() SourceConfig
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

	newBody, ok := sourceBodyFactoryByImpl[b.Impl]
	if !ok {
		var impls []string
		for impl := range sourceBodyFactoryByImpl {
			impls = append(impls, impl)
		}
		sort.Strings(impls)

		return fmt.Errorf(
			"%s: the '%s' source uses '%s' which is not supported, the supported source implementations are: '%s'",
			filename,
			b.Name,
			b.Impl,
			strings.Join(impls, "', '"),
		)
	}

	body := newBody()
	if diag := gohcl.DecodeBody(b.Body, nil, body); diag.HasErrors() {
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

	if err := s.Body.Normalize(cfg); err != nil {
		return fmt.Errorf(
			"%s: the '%s' repository source is invalid: %w",
			s.File,
			s.Block.Name,
			err,
		)
	}

	return nil
}

// assembleSourceBlock converts b into its configuration representation.
func assembleSourceBlock(b sourceBlock, body sourceBlockBody) Source {
	return Source{
		Name:    b.Name,
		Enabled: *b.Enabled,
		Config:  body.Assemble(),
	}
}

var (
	// sourceBodyFactoryByImpl is a map of a source implementation name to a
	// function that returns a new, empty sourceBlockBody type for that
	// implementation.
	sourceBodyFactoryByImpl = map[string]func() sourceBlockBody{}

	// defaultSourceFactoryByName is a map of a source name to a function that
	// returns a new default source. These defaults are merged into any
	// configuration that does not already contain a repository source with the
	// same name.
	defaultSourceFactoryByName = map[string]func() sourceBlockBody{}
)

// registerSourceImpl registers a source implementation, allowing its
// configuration to be parsed.
//
// impl is the name of the implementation, as specified in "source" blocks
// within the configuration file.
func registerSourceImpl(
	impl string,
	newBody func() sourceBlockBody,
) {
	if _, ok := sourceBodyFactoryByImpl[impl]; ok {
		panic("source implementation name already registered")
	}

	sourceBodyFactoryByImpl[impl] = newBody
}

// registerDefaultSource registers a default source that is merged into every
// configuration unless overridden by the user.
func registerDefaultSource(
	name string,
	newBody func() sourceBlockBody,
) {
	if _, ok := defaultSourceFactoryByName[name]; ok {
		panic("default source name already registered")
	}

	defaultSourceFactoryByName[name] = newBody
}
