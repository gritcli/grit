package config

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2"
)

// sourceNameRegexp is a regular expression used to validate source names.
var sourceNameRegexp = regexp.MustCompile(`(?i)^[a-z_]+$`)

// intermediateSource is the intermediate representation of a source.
//
// Some of the contents of a source configuration cannot be finalized until all
// configuration files have been loaded.
type intermediateSource struct {
	Schema sourceSchema
	VCSs   map[string]hcl.Body
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

	i := intermediateSource{
		Schema: s,
		VCSs:   map[string]hcl.Body{},
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
func (l *loader) populateImplicitSources(cfg *Config) error {
	ctx := &sourceContext{
		loader: l,
	}

	for alias, reg := range l.Registry.SourceDrivers() {
		sources, err := reg.ConfigLoader.ImplicitSources(ctx)
		if err != nil {
			return fmt.Errorf(
				"the implicit sources provided by the '%s' driver cannot be loaded: %w",
				alias,
				err,
			)
		}

		for _, src := range sources {
			lowerName := strings.ToLower(src.Name)
			if _, ok := l.sources[lowerName]; ok {
				continue
			}

			cfg.Sources = append(cfg.Sources, Source{
				Name:    src.Name,
				Enabled: true,
				Clones:  l.finalizeImplicitSourceClones(src.Name),
				Driver:  src.Config,
			})
		}
	}

	return nil
}

// finalizeSource returns a source built from its intermediate representation.
func (l *loader) finalizeSource(i intermediateSource) (Source, error) {
	reg, ok := l.Registry.SourceDriverByAlias(i.Schema.Driver)
	if !ok {
		return Source{}, fmt.Errorf(
			"the '%s' source uses an unrecognized driver ('%s'), the supported source drivers are '%s'",
			i.Schema.Name,
			i.Schema.Driver,
			strings.Join(l.Registry.SourceDriverAliases(), "', '"),
		)
	}

	clones, err := l.finalizeSourceSpecificClones(i, i.Schema.Clones)
	if err != nil {
		return Source{}, err
	}

	sourceVCSs, err := l.finalizeSourceSpecificVCSs(i)
	if err != nil {
		return Source{}, err
	}

	ctx := &sourceContext{
		loader:     l,
		sourceVCSs: sourceVCSs,
	}

	cfg, err := reg.ConfigLoader.Defaults(ctx)
	if err != nil {
		if isHCLError(err) {
			return Source{}, err
		}

		return Source{}, fmt.Errorf(
			"the default configuration for the '%s' source driver cannot be loaded: %w",
			i.Schema.Driver,
			err,
		)
	}

	cfg, err = reg.ConfigLoader.Merge(ctx, cfg, i.Schema.DriverBody)
	if err != nil {
		if isHCLError(err) {
			return Source{}, err
		}

		return Source{}, fmt.Errorf(
			"the configuration for the '%s' source cannot be loaded: %w",
			i.Schema.Name,
			err,
		)
	}

	// TODO: produce an error if the source has VCS configurations for
	// unsupported VCS drivers.

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

// sourceContext is an implementation of sourcedriver.ConfigContext.
type sourceContext struct {
	loader     *loader
	sourceVCSs map[string]vcsdriver.Config
}

func (c *sourceContext) EvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{}
}

func (c *sourceContext) NormalizePath(p *string) error {
	return c.loader.normalizePath(p)
}

func (c *sourceContext) UnmarshalVCSConfig(driver string, v interface{}) error {
	configInterfaceType := reflect.TypeOf((*vcsdriver.Config)(nil)).Elem()

	target := reflect.ValueOf(v)

	if v == nil {
		panic(fmt.Sprintf(
			"v must be a pointer to a concrete implementation of the %s interface, but it is nil",
			configInterfaceType,
		))
	}

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

	for alias, reg := range c.loader.Registry.VCSDrivers() {
		if reg.Name != driver {
			continue
		}

		matches = append(matches, alias)

		cfg, ok := c.sourceVCSs[alias]
		if !ok {
			cfg = c.loader.globalVCSs[alias]
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
