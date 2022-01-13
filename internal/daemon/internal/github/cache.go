package github

import (
	"sync"

	"github.com/google/go-github/github"
)

// cache is an in-memory cache of data from the GitHub API.
//
// The values are immutable, they are always replaced completely when set. This
// allows cache consumers to use the values after they have been obtained while
// only requiring the mutex to be held long enough to return the value from the
// cache.
//
// Note that it does not cache all repositories, only those to which the
// authenticated user has been granted explicit access.
type cache struct {
	m            sync.RWMutex
	user         *github.User
	reposByOwner map[string]map[string]*github.Repository
	reposByID    map[int64]*github.Repository
}

// CurrentUser returns the authenticated GitHub user.
func (c *cache) CurrentUser() *github.User {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.user
}

// SetCurrentUser sets the authenticated GitHub user.
func (c *cache) SetCurrentUser(u *github.User) {
	c.m.Lock()
	defer c.m.Unlock()

	c.user = u
}

// ReposByOwner returns map of owner name to repositories.
//
// The second-level map is keyed by the repository name without the owner
// prefix.
func (c *cache) ReposByOwner() map[string]map[string]*github.Repository {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.reposByOwner
}

// RepoByID returns the repository with the given ID.
func (c *cache) RepoByID(id int64) (*github.Repository, bool) {
	c.m.RLock()
	defer c.m.RUnlock()

	r, ok := c.reposByID[id]
	return r, ok
}

// SetRepos updates the cache to include the given repository.
func (c *cache) SetRepos(repos []*github.Repository) {
	reposByOwner := map[string]map[string]*github.Repository{}
	reposByID := map[int64]*github.Repository{}

	for _, r := range repos {
		owner := r.GetOwner()

		reposByName := reposByOwner[owner.GetLogin()]
		if reposByName == nil {
			reposByName = map[string]*github.Repository{}
			reposByOwner[owner.GetLogin()] = reposByName
		}

		reposByName[r.GetName()] = r
		reposByID[r.GetID()] = r
	}

	c.m.Lock()
	defer c.m.Unlock()

	c.reposByOwner = reposByOwner
	c.reposByID = reposByID
}
