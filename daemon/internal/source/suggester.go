package source

import (
	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
)

// A Suggester suggests repositories based on a partial name.
type Suggester struct {
	Sources List
	Log     logs.Log
}

// Suggest suggests a set of repositories that begin with the given word.
func (s *Suggester) Suggest(
	word string,
	includeLocal bool,
	includeRemote bool,
) map[string][]sourcedriver.RemoteRepo {
	// TODO: honour "include" flags

	suggestions := map[string][]sourcedriver.RemoteRepo{}

	for _, src := range s.Sources {
		log := src.
			Log(s.Log).
			WithPrefix("suggest %#v: ", word)

		count := 0
		for w, repos := range src.Driver.Suggest(word, log) {
			count += len(repos)
			suggestions[w] = append(
				suggestions[w],
				repos...,
			)
		}

		log.Write("suggested %d repo(s)", count)
	}

	return suggestions
}
