package github

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/config"
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

	opts, err := newCloneOptions(
		s.cfg,
		r,
		logging.Tee(
			logging.Demote(serverLog), // log to the server as debug
			clientLog,                 // log to the client as regular message
		),
	)
	if err != nil {
		logging.Log(serverLog, "unable to construct clone options: %w", err)
		return err
	}

	_, err = git.PlainCloneContext(
		ctx,
		dir,
		false, // isBare
		opts,
	)

	return err
}

// newCloneOptions returns new clone options based on source configuration.
func newCloneOptions(
	cfg config.GitHub,
	r *github.Repository,
	logger logging.Logger,
) (*git.CloneOptions, error) {
	opts := &git.CloneOptions{
		URL: r.GetSSHURL(),
		Progress: &logging.LineWriter{
			Target: logger,
		},
	}

	if cfg.Git.PreferHTTP {
		opts.URL = r.GetCloneURL()

		if cfg.Token != "" {
			opts.Auth = &http.BasicAuth{
				Username: "token", // ignored
				Password: cfg.Token,
			}
		}
	}

	if cfg.Git.SSHKeyFile != "" {
		publicKeys, err := ssh.NewPublicKeysFromFile(
			"git", // github always uses the "git" username for SSH access
			cfg.Git.SSHKeyFile,
			cfg.Git.SSHKeyPassphrase,
		)
		if err != nil {
			return nil, err
		}

		opts.Auth = publicKeys
	}

	return opts, nil
}
