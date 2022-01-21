package config_test

import (
	"fmt"

	"github.com/gritcli/grit/driver/sourcedriver"
)

// sourceConfigStub is a test implementation of sourcedriver.Config.
type sourceConfigStub struct {
	Value          string
	FilesystemPath string
	VCSConfig      vcsConfigStub
}

// NewDriver constructs a new driver that uses this configuration.
func (c sourceConfigStub) NewDriver() sourcedriver.Driver {
	panic("not implemented")
}

// DescribeSourceConfig returns a human-readable description of the
// configuration.
func (c sourceConfigStub) DescribeSourceConfig() string {
	return fmt.Sprintf(
		"test source (%s)",
		c.Value,
	)
}

// sourceConfigSchemaStub is a test implementation of sourcedriver.ConfigSchema.
type sourceConfigSchemaStub struct {
	Value          string `hcl:"value,optional"`
	FilesystemPath string `hcl:"filesystem_path,optional"`
}

// Normalize validates the configuration as parsed by this schema and
// returns a normalized Config.
func (s *sourceConfigSchemaStub) Normalize(
	ctx sourcedriver.ConfigNormalizeContext,
) (sourcedriver.Config, error) {
	cfg := sourceConfigStub{
		Value:          s.Value,
		FilesystemPath: s.FilesystemPath,
	}

	if cfg.Value == "" {
		cfg.Value = "<default>"
	}

	if err := ctx.NormalizePath(&cfg.FilesystemPath); err != nil {
		return nil, err
	}

	if err := ctx.UnmarshalVCSConfig(testVCSRegistration.Name, &cfg.VCSConfig); err != nil {
		return nil, err
	}

	return cfg, nil
}

// testSourceRegistration contains registration info for the test source driver.
var testSourceRegistration = sourcedriver.Registration{
	Name: "test_source_driver",
	NewConfigSchema: func() sourcedriver.ConfigSchema {
		return &sourceConfigSchemaStub{}
	},
}
