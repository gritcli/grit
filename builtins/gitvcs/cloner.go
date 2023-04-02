package gitvcs

import (
	"context"
	"errors"
	"io"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gritcli/grit/logs"
)

// Cloner is an implementation of sourcedriver.Cloner that clones a Git
// repository.
type Cloner struct {
	// SSHEndpoint is the URL used to clone the repository using the SSH
	// protocol, if available.
	SSHEndpoint string

	// SSHKeyFile is the path to the private SSH key used to authenticate when
	// using the SSH transport.
	//
	// If it is empty, the system's SSH agent is queried to determine which key
	// to use.
	SSHKeyFile string

	// SSHKeyPassphrase is the passphrase used to decrypt the SSH private key,
	// if any. It is ignored if SSHKeyFile is empty.
	SSHKeyPassphrase string

	// HTTPEndpoint is the URL used to clone the repository using the HTTP
	// protocol.
	HTTPEndpoint string

	// HTTPUsername is the username to use when cloning via HTTP, if any.
	HTTPUsername string

	// HTTPPassword is the password to use when cloning via HTTP, if any.
	HTTPPassword string

	// PreferHTTP indicates that the HTTP protocol should be used in preference
	// to SSH. By default SSH is preferred.
	PreferHTTP bool
}

// Clone clones the repository into the given target directory.
func (c *Cloner) Clone(
	ctx context.Context,
	dir string,
	log logs.Log,
) error {
	opts, err := c.cloneOptions(log)
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
func (c *Cloner) cloneOptions(log logs.Log) (*git.CloneOptions, error) {
	useHTTP, err := c.useHTTP()
	if err != nil {
		return nil, err
	}

	if useHTTP {
		return c.httpCloneOptions(log)
	}

	return c.sshCloneOptions(log)
}

// sshCloneOptions returns options that clone the repository using the HTTP
// protocol.
func (c *Cloner) httpCloneOptions(log logs.Log) (*git.CloneOptions, error) {
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
		Progress: progressWriter(log),
	}, nil
}

// sshCloneOptions returns options that clone the repository using the SSH
// protocol.
func (c *Cloner) sshCloneOptions(log logs.Log) (*git.CloneOptions, error) {
	opts := &git.CloneOptions{
		URL:      c.SSHEndpoint,
		Progress: progressWriter(log),
	}

	if c.SSHKeyFile != "" {
		ep, err := transport.NewEndpoint(c.SSHEndpoint)
		if err != nil {
			return nil, err
		}

		publicKeys, err := ssh.NewPublicKeysFromFile(
			ep.User,
			c.SSHKeyFile,
			c.SSHKeyPassphrase,
		)
		if err != nil {
			return nil, err
		}

		opts.Auth = publicKeys
	}

	return opts, nil
}

// useHTTP returns true if the HTTP protocol should be used to clone a
// repository based on a set of configuration values.
//
// The algorithm favours SSH over HTTP unless the configuration specifically
// indicates that HTTP is to be preferred.
//
// If it returns false, SSH should be used instead.
func (c *Cloner) useHTTP() (bool, error) {
	hasSSH := c.SSHEndpoint != ""
	hasHTTP := c.HTTPEndpoint != ""

	if !hasHTTP && !hasSSH {
		return false, errors.New("neither the SSH nor HTTP protocol is available")
	}

	if hasHTTP && c.PreferHTTP {
		return true, nil
	}

	if hasSSH {
		if c.SSHKeyFile != "" {
			return false, nil
		}

		if hasSSHAgent := os.Getenv("SSH_AUTH_SOCK") != ""; hasSSHAgent {
			return false, nil
		}
	}

	if hasHTTP {
		return true, nil
	}

	return false, errors.New("SSH is the only available protocol but there is no SSH agent and no private key was provided")
}

// progressWriter returns the writer used to log the output from Git.
func progressWriter(log logs.Log) io.Writer {
	return &logs.Writer{
		Target: log.WithPrefix("git: "),
	}
}
