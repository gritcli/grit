package config_test

import (
	"fmt"

	"github.com/gritcli/grit/driver/vcsdriver"
)

// vcsConfigStub is a test implementation of vcsdriver.Config.
type vcsConfigStub struct {
	Value          string
	FilesystemPath string
}

// DescribeVCSConfig returns a human-readable description of the
// configuration.
func (c vcsConfigStub) DescribeVCSConfig() string {
	return fmt.Sprintf(
		"test vcs (%s)",
		c.Value,
	)
}

// vcsConfigSchemaStub is a test implementation of vcsdriver.ConfigSchema.
type vcsConfigSchemaStub struct {
	Value          string `hcl:"value,optional"`
	FilesystemPath string `hcl:"filesystem_path,optional"`
}

// NormalizeGlobals validates the global configuration as parsed by this schema,
// and returns a normalized Config.
func (s *vcsConfigSchemaStub) NormalizeGlobals(
	nc vcsdriver.ConfigNormalizeContext,
) (vcsdriver.Config, error) {
	cfg := vcsConfigStub{
		Value:          s.Value,
		FilesystemPath: s.FilesystemPath,
	}

	if cfg.Value == "" {
		cfg.Value = "<default>"
	}

	if err := nc.NormalizePath(&cfg.FilesystemPath); err != nil {
		return nil, err
	}

	return cfg, nil
}

// NormalizeSourceSpecific validates the configuration as parsed by this schema
// within a "source" block and returns a normalized Config.
func (s *vcsConfigSchemaStub) NormalizeSourceSpecific(
	nc vcsdriver.ConfigNormalizeContext,
	g vcsdriver.Config,
) (vcsdriver.Config, error) {
	cfg := g.(vcsConfigStub)

	if s.Value != "" {
		// Note, we concat here (not replace) so that tests can verify that the
		// defaults are available to NormalizeSourceSpecific()
		cfg.Value += s.Value
	}

	if s.FilesystemPath != "" {
		cfg.FilesystemPath = s.FilesystemPath
	}

	if err := nc.NormalizePath(&cfg.FilesystemPath); err != nil {
		return nil, err
	}

	return cfg, nil
}

// testVCSRegistration contains registration info for the test VCS driver.
var testVCSRegistration = vcsdriver.Registration{
	Name: "test_vcs_driver",
	NewConfigSchema: func() vcsdriver.ConfigSchema {
		return &vcsConfigSchemaStub{}
	},
}
