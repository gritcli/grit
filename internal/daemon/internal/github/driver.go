package github

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/daemon/internal/registry"
	"github.com/gritcli/grit/plugin/driver"
)

// impl is an implementation of driver.Driver that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type impl struct {
	config Config
	client *github.Client
	cache  cache
}

// isEnterpriseServer returns true if domain seems to refer to a GitHub
// Enterprise Server installation.
func isEnterpriseServer(domain string) bool {
	return !strings.EqualFold(domain, "github.com")
}

// SourceDriverRegistration returns the registration information for the GitHub
// source driver.
func SourceDriverRegistration() driver.Registration {
	return driver.Registration{
		Name:        "github",
		Description: "Use repositories from GitHub.com or GitHub Enterprise Server.",
		NewConfigSchema: func() driver.ConfigSchema {
			return &configSchema{}
		},
		DefaultSources: map[string]func() driver.ConfigSchema{
			"github": func() driver.ConfigSchema {
				return &configSchema{
					Domain: "github.com",
				}
			},
		},
	}
}

func init() {
	registry.BuiltIns.RegisterSourceDriver(
		"github",
		SourceDriverRegistration(),
	)
}
