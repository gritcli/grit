package githubsource

import (
	"strings"

	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/logs"
)

// Suggest returns a set of repositories that have names beginning with the
// given word (which may be empty).
//
// This implementation considers the name to start with the word if any of the
// owner name, unqualified name, or fully qualified name begin with the word.
func (s *source) Suggest(
	word string,
	log logs.Log,
) map[string][]sourcedriver.RemoteRepo {
	suggestions := map[string][]sourcedriver.RemoteRepo{}

	for _, r := range s.reposByID {
		candidates := []string{
			r.GetFullName(),
			r.GetName(),
		}

		for _, c := range candidates {
			if strings.HasPrefix(c, word) {
				suggestions[c] = append(
					suggestions[c],
					toRemoteRepo(r),
				)
				break
			}
		}
	}

	return suggestions
}
