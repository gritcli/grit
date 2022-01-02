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
	query string,
	out logging.Logger,
) ([]source.Repo, error) {
	out = logging.Tee(
		logging.Prefix(s.logger, "resolve[%s]: ", query),
		out,
	)

	ownerName, repoName, err := parseRepoName(query)
	if err != nil {
		logging.Debug(out, err.Error())
		return nil, nil
	}

	reposByOwner := s.cache.ReposByOwner()
	var repos []source.Repo

	if ownerName == "" {
		for _, reposByName := range reposByOwner {
			if r, ok := reposByName[repoName]; ok {
				repos = append(repos, convertRepo(r))
			}
		}

		logging.Debug(
			out,
			"found %d repo(s) named '%s' by scanning the user's repo cache",
			len(repos),
			repoName,
		)

		return repos, nil
	}

	if r, ok := reposByOwner[ownerName][repoName]; ok {
		logging.Debug(out, "found an exact match in the user's repo cache")

		return []source.Repo{
			convertRepo(r),
		}, nil
	}

	r, res, err := s.client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		if res.StatusCode == http.StatusNotFound {
			logging.Debug(out, "no matches found when querying the API")
			return nil, nil
		}

		return nil, err
	}

	logging.Debug(s.logger, "found an exact match by querying the API")

	return []source.Repo{
		convertRepo(r),
	}, nil
}
