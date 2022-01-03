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

// resolve returns the configuration that is produced by this block.
func (b sourceBlock) resolve(filename string) (Source, error) {
	if b.Name == "" {
		return Source{}, errors.New("source name must not be empty")
	}

	if !sourceNameRegexp.MatchString(b.Name) {
		return Source{}, errors.New("source name must contain only alpha-numeric characters and underscores")
	}

	enabled := true
	if b.Enabled != nil {
		enabled = *b.Enabled
	}

	body, err := decodeSourceBody(b.Impl, b.Body)
	if err != nil {
		return Source{}, err
	}

	cfg, err := body.resolve(filename)
	if err != nil {
		return Source{}, err
	}

	return Source{
		Name:    b.Name,
		Enabled: enabled,
		Config:  cfg,
	}, nil
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
			"'%s' is not recognized source implementation, expected %s",
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
