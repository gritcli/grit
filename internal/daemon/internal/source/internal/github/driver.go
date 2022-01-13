package github

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/daemon/internal/config"
)

// Driver is an implementation of driver.Driver that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type Driver struct {
	Config config.GitHub

	client *github.Client
	cache  cache
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(cfg config.GitHub) bool {
	return strings.EqualFold(cfg.Domain, "github.com")
}
