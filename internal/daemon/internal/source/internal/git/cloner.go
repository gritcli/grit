package git

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gritcli/grit/internal/common/config"
)

// Cloner clones Git repositories.
type Cloner struct {
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

	// Logger is the target for output from the cloning operation.
	Logger logging.Logger
}

// Clone clones the repository into the given target directory.
func (c *Cloner) Clone(ctx context.Context, dir string) error {
	opts, err := c.cloneOptions()
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
func (c *Cloner) cloneOptions() (*git.CloneOptions, error) {
	h, err := useHTTP(
		c.SSHEndpoint != "",
		c.HTTPEndpoint != "",
		c.Config,
	)
	if err != nil {
		return nil, err
	}

	if h {
		return c.httpCloneOptions()
	}

	return c.sshCloneOptions()
}

// sshCloneOptions returns options that clone the repository using the HTTP
// protocol.
func (c *Cloner) httpCloneOptions() (*git.CloneOptions, error) {
	var auth *http.BasicAuth
	if c.HTTPUsername != "" || c.HTTPPassword != "" {
		auth = &http.BasicAuth{
			Username: c.HTTPUsername,
			Password: c.HTTPPassword,
		}
	}

	return &git.CloneOptions{
		URL:  c.HTTPEndpoint,
		Auth: auth,
		Progress: &logging.LineWriter{
			Target: c.Logger,
		},
	}, nil
}

// sshCloneOptions returns options that clone the repository using the SSH
// protocol.
func (c *Cloner) sshCloneOptions() (*git.CloneOptions, error) {
	opts := &git.CloneOptions{
		URL: c.SSHEndpoint,
		Progress: &logging.LineWriter{
			Target: c.Logger,
		},
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
