package config

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// sourceNameRegexp is a regular expression used to validate source names.
var sourceNameRegexp = regexp.MustCompile(`(?i)^[a-z_]+$`)

// prepareSourceBlock prepares b for merging into the configuration once the
// entire configuration has been parsed.
func (l *loader) prepareSourceBlock(filename string, b sourceBlock) error {
	if b.Name == "" {
		return errors.New("repository sources must not have empty names")
	}

	if !sourceNameRegexp.MatchString(b.Name) {
		return fmt.Errorf(
			"the '%s' repository source has an invalid name, names must contain only alpha-numeric characters and underscores",
			b.Name,
		)
	}

	if l.sourceBlockFiles == nil {
		l.sourceBlockFiles = map[string]string{}
	} else if _, ok := l.sourceBlockFiles[b.Name]; ok {
		return fmt.Errorf(
			"a repository source named '%s' has already been defined in %s",
			b.Name,
			l.sourceBlockFiles[b.Name],
		)
	}

	l.sourceBlocks = append(l.sourceBlocks, b)
	l.sourceBlockFiles[b.Name] = filename

	return nil
}

// mergeSourceBlock merges b into the configuration.
//
// It must only be called after the global configuration in l.config has been
// fully parsed, and defaults merged.
func (l *loader) mergeSourceBlock(filename string, b sourceBlock) error {
	src := Source{
		Name:    b.Name,
		Enabled: true,
	}

	if b.Enabled != nil {
		src.Enabled = *b.Enabled
	}

	body, err := decodeSourceBody(b.Impl, b.Body)
	if err != nil {
		return err
	}

	src.Config, err = body.resolve(filename, l.config)
	if err != nil {
		return fmt.Errorf(
			"the '%s' repository source contains invalid configuration: %w",
			b.Name,
			err,
		)
	}

	l.config.Sources[src.Name] = src

	return nil
}

// decodeSourceBody decodes the body of a source block using an
// implementation-specific schema.
func decodeSourceBody(impl string, body hcl.Body) (sourceBlockBody, error) {
	schema, ok := sourceSchemaByImpl[impl]
	if !ok {
		var allowed []string
		for n := range sourceSchemaByImpl {
			allowed = append(allowed, n)
		}
		sort.Strings(allowed)

		var list string
		for i, n := range allowed {
			if i == 0 {
				list += fmt.Sprintf("'%s'", n)
			} else if i == len(allowed)-1 {
				list += fmt.Sprintf(" or '%s'", n)
			} else {
				list += fmt.Sprintf(", '%s'", n)
			}
		}

		return nil, fmt.Errorf(
			"'%s' is not a recognized repository source implementation, expected %s",
			impl,
			list,
		)
	}

	ptr := reflect.New(schema)
	diag := gohcl.DecodeBody(body, nil, ptr.Interface())
	if diag.HasErrors() {
		return nil, diag
	}

	return ptr.Elem().Interface().(sourceBlockBody), nil
}
