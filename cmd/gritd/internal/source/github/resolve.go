package github

import (
	"context"
	"net/http"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
)

// Resolve resolves a repository name, URL, or other identifier to a set of
// possible repositories.
func (d *Driver) Resolve(
	ctx context.Context,
	query string,
	clientLog logging.Logger,
) ([]source.Repo, error) {
	serverLog := logging.Prefix(d.Logger, "resolve[%s]: ", query)
	clientLog = logging.Tee(serverLog, clientLog) // log everything sent to the client on the server as well

	ownerName, repoName, err := parseRepoName(query)
	if err != nil {
		logging.Debug(clientLog, err.Error())
		return nil, nil
	}

	reposByOwner := d.cache.ReposByOwner()
	var repos []source.Repo

	if ownerName == "" {
		for _, reposByName := range reposByOwner {
			if r, ok := reposByName[repoName]; ok {
				repos = append(repos, convertRepo(r))
			}
		}

		logging.Debug(
			clientLog,
			"found %d repo(s) named '%s' by scanning the user's repo cache",
			len(repos),
			repoName,
		)

		return repos, nil
	}

	if r, ok := reposByOwner[ownerName][repoName]; ok {
		logging.Debug(clientLog, "found an exact match in the user's repo cache")

		return []source.Repo{
			convertRepo(r),
		}, nil
	}

	r, res, err := d.client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		if res.StatusCode == http.StatusNotFound {
			logging.Debug(clientLog, "no matches found when querying the API")
			return nil, nil
		}

		logging.Log(serverLog, "unable to query API: %s", err)
		return nil, err
	}

	logging.Debug(serverLog, "found an exact match by querying the API")

	return []source.Repo{
		convertRepo(r),
	}, nil
}
