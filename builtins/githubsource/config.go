package githubsource

import (
	"github.com/gritcli/grit/builtins/gitvcs"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// Config contains configuration specific to the GitHub driver.
type Config struct {
	// Domain is the base domain name of the GitHub installation.
	Domain string

	// Token is a personal access token used to authenticate with the GitHub
	// API.
	Token string

	// Git is the configuration that controls how Grit uses Git for this source.
	Git gitvcs.Config
}

// NewDriver constructs a new driver that uses this configuration.
func (c Config) NewDriver() sourcedriver.Driver {
	return &impl{config: c}
}

// DescribeSourceConfig returns a human-readable description of the
// configuration.
func (c Config) DescribeSourceConfig() string {
	desc := c.Domain

	if isEnterpriseServer(c.Domain) {
		desc += " (github enterprise server)"
	}

	return desc
}

// configSchema is the HCL schema for a "source" block that uses the "github"
// source driver.
type configSchema struct {
	Domain string `hcl:"domain,optional"`
	Token  string `hcl:"token,optional"`
}

// configLoader is an implementation of vcsdriver.ConfigLoader for Git.
type configLoader struct{}

func (configLoader) Defaults(ctx sourcedriver.ConfigContext) (sourcedriver.Config, error) {
	cfg := Config{
		Domain: "github.com",
	}

	if err := ctx.UnmarshalVCSConfig(gitvcs.Registration.Name, &cfg.Git); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (configLoader) Merge(ctx sourcedriver.ConfigContext, c sourcedriver.Config, b hcl.Body) (sourcedriver.Config, error) {
	var s configSchema
	if diag := gohcl.DecodeBody(b, ctx.EvalContext(), &s); diag.HasErrors() {
		return nil, diag
	}

	cfg := c.(Config) // clone

	if s.Domain != "" {
		cfg.Domain = s.Domain
	}

	if s.Token != "" {
		cfg.Token = s.Token
	}

	if err := ctx.UnmarshalVCSConfig(gitvcs.Registration.Name, &cfg.Git); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (l configLoader) ImplicitSources(ctx sourcedriver.ConfigContext) ([]sourcedriver.ImplicitSource, error) {
	cfg, err := l.Defaults(ctx)
	if err != nil {
		return nil, err
	}

	return []sourcedriver.ImplicitSource{
		{
			Name:   "github",
			Config: cfg,
		},
	}, nil
}
