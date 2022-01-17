package fixtures

import (
	"fmt"

	"github.com/gritcli/grit/driver/vcsdriver"
)

// VCSConfigStub is a test implementation of vcsdriver.Config.
type VCSConfigStub struct {
	Value string
}

// DescribeVCSConfig returns a human-readable description of the
// configuration.
func (c VCSConfigStub) DescribeVCSConfig() string {
	return fmt.Sprintf(
		"test vcs (%s)",
		c.Value,
	)
}

// VCSConfigSchemaStub is a test implementation of vcsdriver.ConfigSchema.
type VCSConfigSchemaStub struct {
	Value string `hcl:"value,optional"`
}

// NormalizeDefaults validates the configuration as parsed by this schema at the
// "top-level" of a Grit configuration, and returns a normalized Config.
func (s *VCSConfigSchemaStub) NormalizeDefaults(
	ctx vcsdriver.ConfigNormalizeContext,
) (vcsdriver.Config, error) {
	cfg := VCSConfigStub{
		Value: s.Value,
	}

	if cfg.Value == "" {
		cfg.Value = "<default>"
	}

	return cfg, nil
}

// NormalizeSourceSpecific validates the configuration as parsed by this schema
// within a "source" block and returns a normalized Config.
func (s *VCSConfigSchemaStub) NormalizeSourceSpecific(
	ctx vcsdriver.ConfigNormalizeContext,
	def vcsdriver.Config,
) (vcsdriver.Config, error) {
	cfg := def.(VCSConfigStub)

	if s.Value != "" {
		cfg.Value += s.Value
	}

	return cfg, nil
}

// VCSRegistration contains registration info for the test VCS driver.
var VCSRegistration = vcsdriver.Registration{
	Name:        "test_vcs_driver",
	Description: "test VCS driver",
	NewConfigSchema: func() vcsdriver.ConfigSchema {
		return &VCSConfigSchemaStub{}
	},
}
