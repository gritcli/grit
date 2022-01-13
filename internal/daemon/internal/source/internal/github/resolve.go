package github

import (
	"context"
	"net/http"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/plugin/driver"
)

// Resolve resolves a repository name, URL, or other identifier to a set of
// possible repositories.
func (d *Driver) Resolve(
	ctx context.Context,
	query string,
	logger logging.Logger,
) ([]driver.RemoteRepo, error) {
	ownerName, repoName, err := parseRepoName(query)
	if err != nil {
		return nil, nil
	}

	reposByOwner := d.cache.ReposByOwner()
	var repos []*github.Repository

	if ownerName == "" {
		for _, reposByName := range reposByOwner {
			if r, ok := reposByName[repoName]; ok {
				repos = append(repos, r)
			}
		}

		logging.Debug(
			logger,
			"found %d match(es) for '%s' in the repository list for @%s",
			len(repos),
			query,
			d.cache.CurrentUser().GetLogin(),
		)

		if len(repos) == 0 {
			logging.Debug(
				logger,
				"skipping GitHub API query for '%s' because it is not a fully-qualified repository name",
				query,
			)
		}

		return toRemoteRepos(repos...), nil
	}

	if r, ok := reposByOwner[ownerName][repoName]; ok {
		logging.Debug(
			logger,
			"found an exact match for '%s' in the repository list for @%s",
			query,
			d.cache.CurrentUser().GetLogin(),
		)

		return toRemoteRepos(r), nil
	}

	r, res, err := d.client.Repositories.Get(ctx, ownerName, repoName)
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
