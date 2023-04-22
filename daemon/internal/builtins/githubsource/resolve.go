package githubsource

import (
	"context"
	"net/http"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
)

// Resolve resolves a repository name, URL, or other identifier to a set of
// possible repositories.
func (s *source) Resolve(
	ctx context.Context,
	query string,
	log logs.Log,
) ([]sourcedriver.RemoteRepo, error) {
	ownerName, repoName, err := parseRepoName(query)
	if err != nil {
		return nil, nil
	}

	state := s.state.Load()

	if ownerName == "" {
		var matches []sourcedriver.RemoteRepo

		for _, reposByName := range state.ReposByOwner {
			if r, ok := reposByName[repoName]; ok {
				matches = append(matches, toRemoteRepo(r))
			}
		}

		log.WriteVerbose(
			"found %d match(es) for '%s' in the repository list for @%s",
			len(matches),
			query,
			state.User.GetLogin(),
		)

		if len(matches) == 0 {
			log.WriteVerbose(
				"skipping GitHub API query for '%s' because it is not a fully-qualified repository name",
				query,
			)
		}

		return matches, nil
	}

	if r, ok := state.ReposByOwner[ownerName][repoName]; ok {
		log.WriteVerbose(
			"found an exact match for '%s' in the repository list for @%s",
			query,
			state.User.GetLogin(),
		)

		return toRemoteRepos(r), nil
	}

	r, res, err := state.Client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		if res != nil || res.StatusCode == http.StatusNotFound {
			log.WriteVerbose(
				"no repository named '%s' found by querying the GitHub API",
				query,
			)

			return nil, nil
		}

		return nil, err
	}

	log.WriteVerbose(
		"found a repository named '%s' by querying the GitHub API",
		query,
	)

	return toRemoteRepos(r), nil
}
