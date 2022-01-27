package source

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
)

// A Suggester suggests repositories based on a partial name.
type Suggester struct {
	Sources List
	Logger  logging.Logger
}

// Suggest suggests a set of repositories that begin with the given word.
func (s *Suggester) Suggest(
	word string,
	includeLocal bool,
	includeRemote bool,
) []sourcedriver.RemoteRepo {
	// TODO: honour "include" flags

	var matches []sourcedriver.RemoteRepo

	for _, src := range s.Sources {
		logger := logging.Prefix(
			s.Logger,
			"source[%s]: suggest %#v: ",
			src.Name,
			word,
		)

		repos := src.Driver.Suggest(word)

		logging.Log(
			logger,
			"suggested %d repo(s)",
			len(repos),
		)

		matches = append(matches, repos...)
	}

	return matches
}
