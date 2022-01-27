package githubsource

import (
	"context"
	"net/http"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
)

// Resolve resolves a repository name, URL, or other identifier to a set of
// possible repositories.
func (s *source) Resolve(
	ctx context.Context,
	query string,
	logger logging.Logger,
) ([]sourcedriver.RemoteRepo, error) {
	ownerName, repoName, err := parseRepoName(query)
	if err != nil {
		return nil, nil
	}

	if ownerName == "" {
		var matches []sourcedriver.RemoteRepo

		for _, reposByName := range s.reposByOwner {
			if r, ok := reposByName[repoName]; ok {
				matches = append(matches, toRemoteRepo(r))
			}
		}

		logging.Debug(
			logger,
			"found %d match(es) for '%s' in the repository list for @%s",
			len(matches),
			query,
			s.user.GetLogin(),
		)

		if len(matches) == 0 {
			logging.Debug(
				logger,
				"skipping GitHub API query for '%s' because it is not a fully-qualified repository name",
				query,
			)
		}

		return matches, nil
	}

	if r, ok := s.reposByOwner[ownerName][repoName]; ok {
		logging.Debug(
			logger,
			"found an exact match for '%s' in the repository list for @%s",
			query,
			s.user.GetLogin(),
		)

		return toRemoteRepos(r), nil
	}

	r, res, err := s.client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		if res.StatusCode == http.StatusNotFound {
			logging.Debug(
				logger,
				"no repository named '%s' found by querying the GitHub API",
				query,
			)

			return nil, nil
		}

		return nil, err
	}

	logging.Debug(
		logger,
		"found a repository named '%s' by querying the GitHub API",
		query,
	)

	return toRemoteRepos(r), nil
}
