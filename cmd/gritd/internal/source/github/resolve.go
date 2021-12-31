package github

import (
	"context"
	"net/http"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
)

// Resolve resolves a repository name, URL, or other identifier to a set of
// possible repositories.
func (s *impl) Resolve(
	ctx context.Context,
	name string,
	out logging.Logger,
) ([]source.Repo, error) {
	ownerName, repoName, err := parseRepoName(name)
	if err != nil {
		logging.Debug(
			s.logger,
			"resolve[%s]: repository name is not valid: %s",
			name,
			err,
		)

		return nil, nil
	}

	s.repoCacheM.RLock()
	reposByOwner := s.reposByOwner
	s.repoCacheM.RUnlock()

	var results []source.Repo

	if ownerName == "" {
		for _, reposByName := range reposByOwner {
			if r, ok := reposByName[repoName]; ok {
				results = append(results, convertRepo(r))
			}
		}

		logging.Debug(
			s.logger,
			"resolve[%s]: owner not known, found %d repo(s) named '%s' by scanning the user's repo cache",
			name,
			len(results),
			repoName,
		)

		return results, nil
	}

	if r, ok := reposByOwner[ownerName][repoName]; ok {
		logging.Debug(
			s.logger,
			"resolve[%s]: found an exact match in the user's repo cache",
			name,
			len(results),
		)

		return []source.Repo{convertRepo(r)}, nil
	}

	r, res, err := s.client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		if res.StatusCode == http.StatusNotFound {
			logging.Debug(
				s.logger,
				"resolve[%s]: no matches found when querying the API",
				name,
				len(results),
			)

			return nil, nil
		}

		return nil, err
	}

	logging.Debug(
		s.logger,
		"resolve[%s]: found an exact match by querying the API",
		name,
		len(results),
	)

	return []source.Repo{
		convertRepo(r),
	}, nil
}
