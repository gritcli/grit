package githubsource

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/gritcli/grit/driver/sourcedriver"
)

// Suggest returns a set of repositories with names that begin with the
// given word.
func (s *source) Suggest(word string) []sourcedriver.RemoteRepo {
	var matches []*github.Repository

	for owner, repos := range s.reposByOwner {
		if strings.HasPrefix(owner, word) {
			for _, r := range repos {
				matches = append(matches, r)
			}
		} else {
			for name, r := range repos {
				if strings.HasPrefix(r.GetFullName(), word) ||
					strings.HasPrefix(name, word) {
					matches = append(matches, r)
				}
			}
		}
	}

	return toRemoteRepos(matches...)
}
