package config

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
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
func (l *loader) mergeSource(file string, s sourceSchema) error {
	if s.Name == "" {
		return errors.New("source configurations must provide a name")
	}

	if !sourceNameRegexp.MatchString(s.Name) {
		return fmt.Errorf(
			"'%s' is not a valid source name, valid characters are ASCII letters, numbers and underscore",
			s.Name,
		)
	}

	lowerName := strings.ToLower(s.Name)

	if existingSource, ok := l.sources[lowerName]; ok {
		return fmt.Errorf(
			"the '%s' source conflicts with a source of the same name in %s (source names are case-insensitive)",
			s.Name,
			existingSource.File,
		)
	}

	if s.Driver == "" {
		return fmt.Errorf(
			"the '%s' source has an empty driver name",
			s.Name,
		)
	}

	reg, ok := l.Registry.SourceDriverByAlias(s.Driver)
	if !ok {
		return fmt.Errorf(
			"the '%s' source uses an unrecognized driver ('%s'), the supported source drivers are '%s'",
			s.Name,
			s.Driver,
			strings.Join(l.Registry.SourceDriverAliases(), "', '"),
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
		if err := l.mergeSourceSpecificVCS(&i, v); err != nil {
			return err
		}
	}

	if l.sources == nil {
		l.sources = map[string]intermediateSource{}
	}

	l.sources[lowerName] = i

	return nil
}

// populateImplicitSources adds any implicit sources provided by the supported
// source drivers.
//
// Any implicit source with a name that has already been defined in the
// configuration files is ignored.
func (l *loader) populateImplicitSources() {
	for _, reg := range l.Registry.SourceDrivers() {
		for name, newSchema := range reg.ImplicitSources {
			lowerName := strings.ToLower(name)

			if _, ok := l.sources[lowerName]; ok {
				continue
			}

			l.sources[lowerName] = intermediateSource{
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
func (l *loader) finalizeSource(i intermediateSource) (Source, error) {
	clones, err := l.finalizeSourceSpecificClones(i, i.Schema.Clones)
	if err != nil {
		return Source{}, err
	}

	sourceVCSs, err := l.finalizeSourceSpecificVCSs(i)
	if err != nil {
		return Source{}, err
	}

	nc := &sourceNormalizeContext{
		loader:     l,
		globalVCSs: l.globalVCSs,
		sourceVCSs: sourceVCSs,
	}

	cfg, err := i.Driver.Normalize(nc)
	if err != nil {
		return Source{}, fmt.Errorf(
			"the '%s' source is invalid: %w", // TODO: not necessarily invalid config, but just can't be loaded
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
	loader     *loader
	globalVCSs map[string]vcsdriver.Config
	sourceVCSs map[string]vcsdriver.Config
}

func (nc *sourceNormalizeContext) NormalizePath(p *string) error {
	return nc.loader.normalizePath(p)
}

func (nc *sourceNormalizeContext) UnmarshalVCSConfig(driver string, v interface{}) error {
	configInterfaceType := reflect.TypeOf((*vcsdriver.Config)(nil)).Elem()

	target := reflect.ValueOf(v)

	if target.Kind() != reflect.Ptr {
		panic(fmt.Sprintf(
			"v must be a pointer to a concrete implementation of the %s interface, but %s is not a pointer",
			configInterfaceType,
			target.Type(),
		))
	}

	target = target.Elem()

	if !target.Type().Implements(configInterfaceType) {
		panic(fmt.Sprintf(
			"v must be a pointer to a concrete implementation of the %s interface, but %s does not implement that interface",
			configInterfaceType,
			target.Type(),
		))
	}

	if target.Kind() == reflect.Interface {
		panic(fmt.Sprintf(
			"v must be a pointer to a concrete implementation of the %s interface, but %s is not a concrete type (it's an interface)",
			configInterfaceType,
			target.Type(),
		))
	}

	var matches []string

	for alias, reg := range nc.loader.Registry.VCSDrivers() {
		if reg.Name != driver {
			continue
		}

		matches = append(matches, alias)

		cfg, ok := nc.sourceVCSs[alias]
		if !ok {
			cfg = nc.globalVCSs[alias]
		}

		rv := reflect.ValueOf(cfg)

		if rv.Type() == target.Type() {
			target.Set(rv)
			return nil
		}
	}

	if len(matches) == 0 {
		return fmt.Errorf(
			"dependency on unrecognized version control system ('%s')",
			driver,
		)
	}

	sort.Strings(matches)

	return fmt.Errorf(
		"depends on incompatible version control system ('%s'), none of the matching drivers ('%s') use the same configuration structure",
		driver,
		strings.Join(matches, "', '"),
	)
}
