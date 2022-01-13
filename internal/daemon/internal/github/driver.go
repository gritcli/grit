package github

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/daemon/internal/config"
	"github.com/gritcli/grit/plugin/driver"
)

// impl is an implementation of driver.Driver that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type impl struct {
	config config.GitHub
	client *github.Client
	cache  cache
}

// NewDriver returns a new GitHub driver.
func NewDriver(cfg config.GitHub) driver.Driver {
	return &impl{config: cfg}
}

// isEnterpriseServer returns true if domain seems to refer to a GitHub
// Enterprise Server installation.
func isEnterpriseServer(domain string) bool {
	return !strings.EqualFold(domain, "github.com")
}
