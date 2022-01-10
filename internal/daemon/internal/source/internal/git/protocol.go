package git

import (
	"errors"
	"os"

	"github.com/gritcli/grit/internal/common/config"
)

// useHTTP returns true if the HTTP protocol should be used to clone a
// repository based on a set of configuration values.
//
// The algorithm favours SSH over HTTP unless the configuration specifically
// indicates that HTTP is to be preferred.
//
// If it returns false, SSH should be used instead.
func useHTTP(
	hasSSH, hasHTTP bool,
	cfg config.Git,
) (bool, error) {
	if !hasHTTP && !hasSSH {
		return false, errors.New("neither the SSH nor HTTP protocol is available")
	}

	if hasHTTP && cfg.PreferHTTP {
		return true, nil
	}

	if hasSSH {
		if cfg.SSHKeyFile != "" {
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
