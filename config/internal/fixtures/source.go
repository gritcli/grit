package fixtures

import (
	"fmt"

	"github.com/gritcli/grit/driver/sourcedriver"
)

// SourceConfigStub is a test implementation of sourcedriver.Config.
type SourceConfigStub struct {
	Value     string
	VCSConfig VCSConfigStub
}

// NewDriver constructs a new driver that uses this configuration.
func (c SourceConfigStub) NewDriver() sourcedriver.Driver {
	panic("not implemented")
}

// DescribeSourceConfig returns a human-readable description of the
// configuration.
func (c SourceConfigStub) DescribeSourceConfig() string {
	return fmt.Sprintf(
		"test source (%s)",
		c.Value,
	)
}

// SourceConfigSchemaStub is a test implementation of sourcedriver.ConfigSchema.
type SourceConfigSchemaStub struct {
	Value string `hcl:"value,optional"`
}

// Normalize validates the configuration as parsed by this schema and
// returns a normalized Config.
func (s *SourceConfigSchemaStub) Normalize(
	ctx sourcedriver.ConfigNormalizeContext,
) (sourcedriver.Config, error) {
	cfg := SourceConfigStub{
		Value: s.Value,
	}

	if cfg.Value == "" {
		cfg.Value = "<default>"
	}

	if err := ctx.ResolveVCSConfig(&cfg.VCSConfig); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SourceRegistration contains registration info for the test source driver.
var SourceRegistration = sourcedriver.Registration{
	Name:        "test_source_driver",
	Description: "test source driver",
	NewConfigSchema: func() sourcedriver.ConfigSchema {
		return &SourceConfigSchemaStub{}
	},
}
