package github

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
)

// Clone makes a repository available at the specified directory.
func (s *impl) Clone(
	ctx context.Context,
	repoID, dir string,
	clientLog logging.Logger,
) error {
	serverLog := logging.Prefix(s.logger, "clone[%s]: ", repoID)

	id, err := parseRepoID(repoID)
	if err != nil {
		logging.LogString(serverLog, err.Error())
		return err
	}

	r, ok := s.cache.RepoByID(id)
	if !ok {
		var err error
		r, _, err = s.client.Repositories.GetByID(ctx, id)
		if err != nil {
			logging.Log(serverLog, "unable to query API: %s", err)
			return err
		}
	}

	logging.Debug(serverLog, "cloning %s to %s", r.GetFullName(), dir)

	_, err = git.PlainCloneContext(
		ctx,
		dir,
		false, // isBare
		&git.CloneOptions{
			URL: r.GetSSHURL(),
			Progress: &logging.LineWriter{
				Target: logging.Tee(
					logging.Demote(serverLog), // log to the server as debug
					clientLog,                 // log to the client as regular message
				),
			},
		},
	)

	return err
}
