package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// sourceNameRegexp is a regular expression used to validate source names.
var sourceNameRegexp = regexp.MustCompile(`(?i)^[a-z_]+$`)

// intermediateSource is the intermediate representation of a source.
//
// Some of the contents of a source configuration can not be finalized until all
// configuration files have been loaded.
type intermediateSource struct {
	Schema sourceSchema
	Driver sourcedriver.ConfigSchema
	VCSs   map[string]vcsdriver.ConfigSchema
	File   string
}

// mergeSource merges s into the configuration.
func (r *resolver) mergeSource(file string, s sourceSchema) error {
	if s.Name == "" {
		return fmt.Errorf(
			"%s: source configurations must provide a name",
			file,
		)
	}

	if !sourceNameRegexp.MatchString(s.Name) {
		return fmt.Errorf(
			"%s: '%s' is not a valid source name, valid characters are ASCII letters, numbers and underscore",
			file,
			s.Name,
		)
	}

	lowerName := strings.ToLower(s.Name)

	if existingSource, ok := r.sources[lowerName]; ok {
		return fmt.Errorf(
			"%s: the '%s' source conflicts with a source of the same name in %s (source names are case-insensitive)",
			file,
			s.Name,
			existingSource.File,
		)
	}

	if s.Driver == "" {
		return fmt.Errorf(
			"%s: the '%s' source has an empty driver name",
			file,
			s.Name,
		)
	}

	reg, ok := r.registry.SourceDriverByAlias(s.Driver)
	if !ok {
		return fmt.Errorf(
			"%s: the '%s' source uses an unrecognized driver ('%s'), the supported source drivers are '%s'",
			file,
			s.Name,
			s.Driver,
			strings.Join(r.registry.SourceDriverAliases(), "', '"),
		)
	}

	bodySchema := reg.NewConfigSchema()
	if diag := gohcl.DecodeBody(s.DriverBody, nil, bodySchema); diag.HasErrors() {
		return diag
	}

	i := intermediateSource{
		Schema: s,
		Driver: bodySchema,
		VCSs:   map[string]vcsdriver.ConfigSchema{},
		File:   file,
	}

	for _, v := range s.VCSs {
		if err := r.mergeSourceSpecificVCS(&i, v); err != nil {
			return err
		}
	}

	if r.sources == nil {
		r.sources = map[string]intermediateSource{}
	}

	r.sources[lowerName] = i

	return nil
}

// populateImplicitSources adds any implicit sources provided by the supported
// source drivers.
//
// Any implicit source with a name that has already been defined in the
// configuration files is ignored.
func (r *resolver) populateImplicitSources() {
	for _, reg := range r.registry.SourceDrivers() {
		for name, newSchema := range reg.ImplicitSources {
			lowerName := strings.ToLower(name)

			if _, ok := r.sources[lowerName]; ok {
				continue
			}

			r.sources[lowerName] = intermediateSource{
				Schema: sourceSchema{
					Name:   name,
					Driver: reg.Name,
				},
				Driver: newSchema(),
			}
		}
	}
}

// finalizeSource returns a source built from its intermediate representation.
func (r *resolver) finalizeSource(i intermediateSource) (Source, error) {
	clones, err := r.finalizeSourceSpecificClones(i, i.Schema.Clones)
	if err != nil {
		return Source{}, err
	}

	sourceVCSs, err := r.finalizeSourceSpecificVCSs(i)
	if err != nil {
		return Source{}, err
	}

	nc := &sourceNormalizeContext{
		resolver:   r,
		globalVCSs: r.globalVCSs,
		sourceVCSs: sourceVCSs,
	}

	cfg, err := i.Driver.Normalize(nc)
	if err != nil {
		return Source{}, fmt.Errorf(
			"%s: the '%s' source is invalid: %w",
			i.File, // TODO: s.File is empty for "implicit" sources
			i.Schema.Name,
			err,
		)
	}

	enabled := true
	if i.Schema.Enabled != nil {
		enabled = *i.Schema.Enabled
	}

	return Source{
		Name:    i.Schema.Name,
		Enabled: enabled,
		Clones:  clones,
		Driver:  cfg,
	}, nil
}

// sourceNormalizeContext is an implementation of
// sourcedriver.ConfigNormalizeContext.
type sourceNormalizeContext struct {
	resolver   *resolver
	globalVCSs map[string]vcsdriver.Config
	sourceVCSs map[string]vcsdriver.Config
}

func (nc *sourceNormalizeContext) NormalizePath(p *string) error {
	return nc.resolver.normalizePath(p)
}

func (nc *sourceNormalizeContext) ResolveVCSConfig(cfg interface{}) error {
	// TODO: validate that cfg is a pointer-to-impl to provide a better error
	// message.
	elem := reflect.ValueOf(cfg).Elem()

	if nc.resolveVCSConfig(nc.sourceVCSs, elem) {
		return nil
	}

	if nc.resolveVCSConfig(nc.globalVCSs, elem) {
		return nil
	}

	return fmt.Errorf(
		"none of the supported VCS drivers provided a config of type %s",
		elem.Type(),
	)
}

// resolveVCSConfig assigns a value from configs to elem, if possible.
func (nc *sourceNormalizeContext) resolveVCSConfig(
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
