package githubsource

import (
	"strings"

	"github.com/gritcli/grit/driver/sourcedriver"
)

// Suggest returns a set of repositories that have names beginning with the
// given word (which may be empty).
//
// This implementation considers the name to start with the word if any of the
// owner name, unqualified name, or fully qualified name begin with the word.
func (s *source) Suggest(word string) []sourcedriver.RemoteRepo {
	var matches []sourcedriver.RemoteRepo

	for owner, repos := range s.reposByOwner {
		if strings.HasPrefix(owner, word) {
			for _, r := range repos {
				matches = append(matches, toRemoteRepo(r))
			}
		} else {
			for name, r := range repos {
				if strings.HasPrefix(r.GetFullName(), word) ||
					strings.HasPrefix(name, word) {
					matches = append(matches, toRemoteRepo(r))
				}
			}
		}
	}

	return matches
}
