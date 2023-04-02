package githubsource

import (
	"strings"

	"github.com/google/go-github/v50/github"
)

// source is an implementation of sourcedriver.Source that provides repositories
// from GitHub.com or a GitHub Enterprise Server installation.
type source struct {
	config Config
	client *github.Client

	user         *github.User
	reposByID    map[int64]*github.Repository
	reposByOwner map[string]map[string]*github.Repository
}

// isEnterpriseServer returns true if domain seems to refer to a GitHub
// Enterprise Server installation.
func isEnterpriseServer(domain string) bool {
	return !strings.EqualFold(domain, "github.com")
}
