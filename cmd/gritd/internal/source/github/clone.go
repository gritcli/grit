package github

import (
	"context"
	"os"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/common/config"
)

// Clone makes a repository available at the specified directory.
func (d *Driver) Clone(
	ctx context.Context,
	repoID, tempDir string,
	clientLog logging.Logger,
) (string, error) {
	serverLog := logging.Prefix(d.Logger, "clone[%s]: ", repoID)

	id, err := parseRepoID(repoID)
	if err != nil {
		logging.LogString(serverLog, err.Error())
		return "", err
	}

	r, ok := d.cache.RepoByID(id)
	if !ok {
		var err error
		r, _, err = d.client.Repositories.GetByID(ctx, id)
		if err != nil {
			logging.Log(serverLog, "unable to query API: %s", err)
			return "", err
		}
	}

	logging.Debug(serverLog, "cloning %s to %s", r.GetFullName(), tempDir)

	opts, err := newCloneOptions(
		d.Config,
		r,
		logging.Tee(
			logging.Demote(serverLog), // log to the server as debug
			clientLog,                 // log to the client as regular message
		),
	)
	if err != nil {
		logging.Log(serverLog, "unable to construct clone options: %w", err)
		return "", err
	}

	_, err = git.PlainCloneContext(
		ctx,
		tempDir,
		false, // isBare
		opts,
	)

	return r.GetFullName(), err
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
		// The user explicitly prefers HTTP.
		useHTTP(opts, cfg, r)
	} else if cfg.Git.SSHKeyFile != "" {
		// The user prefers SSH and has supplied a specific private key.
		if err := useSSHKey(opts, cfg); err != nil {
			return opts, err
		}
	} else if os.Getenv("SSH_AUTH_SOCK") == "" {
		// The user prefers SSH, but did not provide a private key, and the SSH
		// agent is unavailable, so fall back to HTTP.
		useHTTP(opts, cfg, r)
	}

	return opts, nil
}

// useHTTP configures opts to clone using the HTTP protocol (instead of SSH).
func useHTTP(opts *git.CloneOptions, cfg config.GitHub, r *github.Repository) {
	opts.URL = r.GetCloneURL()

	if cfg.Token != "" {
		opts.Auth = &http.BasicAuth{
			Username: "token", // ignored
			Password: cfg.Token,
		}
	}
}

// useSSHKey configures opts to use the given SSH key for authentication.
func useSSHKey(opts *git.CloneOptions, cfg config.GitHub) error {
	publicKeys, err := ssh.NewPublicKeysFromFile(
		"git", // github always uses the "git" username for SSH access
		cfg.Git.SSHKeyFile,
		cfg.Git.SSHKeyPassphrase,
	)
	if err != nil {
		return err
	}

	opts.Auth = publicKeys

	return nil
}
