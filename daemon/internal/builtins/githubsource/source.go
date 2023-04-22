package githubsource

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/google/go-github/v50/github"
)

// source is an implementation of sourcedriver.Source that provides repositories
// from GitHub.com or a GitHub Enterprise Server installation.
type source struct {
	config Config
	reinit chan struct{}
	state  atomic.Pointer[state]
}

// isEnterpriseServer returns true if domain seems to refer to a GitHub
// Enterprise Server installation.
func isEnterpriseServer(domain string) bool {
	return !strings.EqualFold(domain, "github.com")
}

// newClient returns a new GitHub client for the given configuration.
func newClient(c Config, h *http.Client) (*github.Client, error) {
	if isEnterpriseServer(c.Domain) {
		return github.NewEnterpriseClient(c.Domain, "", h)
	}

	return github.NewClient(h), nil
}

type state struct {
	Client       *github.Client
	User         *github.User
	ReposByID    map[int64]*github.Repository
	ReposByOwner map[string]map[string]*github.Repository
}
