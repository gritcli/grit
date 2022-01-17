package github

import (
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/internal/daemon/internal/config"
)

// Config contains configuration specific to the GitHub driver.
type Config struct {
	// Domain is the base domain name of the GitHub installation.
	Domain string

	// Token is a personal access token used to authenticate with the GitHub
	// API.
	Token string

	// Git is the configuration that controls how Grit uses Git for this source.
	Git config.Git
}

// NewDriver constructs a new driver that uses this configuration.
func (c Config) NewDriver() sourcedriver.Driver {
	return &impl{config: c}
}

// String returns a human-readable description of the configuration.
func (c Config) String() string {
	desc := c.Domain

	if isEnterpriseServer(c.Domain) {
		desc += " (github enterprise server)"
	}

	return desc
}

// configSchema is the HCL schema for a "source" block that uses the "github"
// source driver.
type configSchema struct {
	Domain string                   `hcl:"domain,optional"`
	Token  string                   `hcl:"token,optional"`
	Git    *config.GitBlockExported `hcl:"git,block"`
}

func (s *configSchema) Resolve(ctx sourcedriver.ConfigResolutionContext) (sourcedriver.Config, error) {
	if s.Domain == "" {
		s.Domain = "github.com"
	}

	cfg := Config{
		Domain: s.Domain,
		Token:  s.Token,
	}

	if err := ctx.ResolveVCSConfig(s.Git, &cfg.Git); err != nil {
		return nil, err
	}

	return cfg, nil
}
