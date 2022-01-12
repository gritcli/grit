package git

import (
	"context"
	"io"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gritcli/grit/internal/daemon/internal/config"
)

// BoundCloner is an implementation of source.BoundCloner that clones a Git
// repository.
type BoundCloner struct {
	// Config is the Git configuration for the source that returned the cloner.
	Config config.Git

	// SSHEndpoint is the URL used to clone the repository using the SSH
	// protocol.
	SSHEndpoint string

	// HTTPEndpoint is the URL used to clone the repository using the HTTP
	// protocol.
	HTTPEndpoint string

	// HTTPUsername is the username to use when cloning via HTTP, if any.
	HTTPUsername string

	// HTTPPassword is the password to use when cloning via HTTP, if any.
	HTTPPassword string
}

// Clone clones the repository into the given target directory.
func (c *BoundCloner) Clone(
	ctx context.Context,
	dir string,
	logger logging.Logger,
) error {
	opts, err := c.cloneOptions(logger)
	if err != nil {
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

// cloneOptions returns the options to use when cloning the repository, based on
// the configuration of the cloner.
func (c *BoundCloner) cloneOptions(logger logging.Logger) (*git.CloneOptions, error) {
	h, err := useHTTP(
		c.SSHEndpoint != "",
		c.HTTPEndpoint != "",
		c.Config,
	)
	if err != nil {
		return nil, err
	}

	if h {
		return c.httpCloneOptions(logger)
	}

	return c.sshCloneOptions(logger)
}

// sshCloneOptions returns options that clone the repository using the HTTP
// protocol.
func (c *BoundCloner) httpCloneOptions(logger logging.Logger) (*git.CloneOptions, error) {
	var auth *http.BasicAuth
	if c.HTTPUsername != "" || c.HTTPPassword != "" {
		auth = &http.BasicAuth{
			Username: c.HTTPUsername,
			Password: c.HTTPPassword,
		}
	}

	return &git.CloneOptions{
		URL:      c.HTTPEndpoint,
		Auth:     auth,
		Progress: progressWriter(logger),
	}, nil
}

// sshCloneOptions returns options that clone the repository using the SSH
// protocol.
func (c *BoundCloner) sshCloneOptions(logger logging.Logger) (*git.CloneOptions, error) {
	opts := &git.CloneOptions{
		URL:      c.SSHEndpoint,
		Progress: progressWriter(logger),
	}

	if c.Config.SSHKeyFile != "" {
		ep, err := transport.NewEndpoint(c.SSHEndpoint)
		if err != nil {
			return nil, err
		}

		publicKeys, err := ssh.NewPublicKeysFromFile(
			ep.User,
			c.Config.SSHKeyFile,
			c.Config.SSHKeyPassphrase,
		)
		if err != nil {
			return nil, err
		}

		opts.Auth = publicKeys
	}

	return opts, nil
}

// progressWriter returns the writer used to log the output from Git.
func progressWriter(logger logging.Logger) io.Writer {
	return &logging.StreamWriter{
		Target: logging.Prefix(logger, "git: "),
	}
}
