package gitvcs

import (
	"path/filepath"

	"github.com/gritcli/grit/driver/vcsdriver"
)

// Config is the configuration that controls how Grit uses the Git VCS.
type Config struct {
	// SSHKeyFile is the path to the private SSH key used to authenticate when
	// using the SSH transport.
	//
	// If it is empty, the system's SSH agent is queried to determine which
	// identity to use.
	SSHKeyFile string

	// SSHKeyPassphrase is the passphrase used to encrypt the SSH private key,
	// if any.
	SSHKeyPassphrase string

	// PreferHTTP indicates that Grit should prefer the HTTP transport. By
	// default SSH is preferred.
	PreferHTTP bool
}

// DescribeVCSConfig returns a human-readable description of the
// configuration.
func (c Config) DescribeVCSConfig() string {
	var desc string

	if c.SSHKeyFile != "" {
		desc += "use ssh key (" + filepath.Base(c.SSHKeyFile) + ")"
	} else {
		desc += "use ssh agent"
	}

	if c.PreferHTTP {
		desc += ", prefer http"
	}

	return desc
}

type configSchema struct {
	SSHKey *struct {
		File       string `hcl:"file"`
		Passphrase string `hcl:"passphrase,optional"`
	} `hcl:"ssh_key,block"`
	PreferHTTP *bool `hcl:"prefer_http"`
}

func (s *configSchema) NormalizeDefaults(
	ctx vcsdriver.ConfigNormalizeContext,
) (vcsdriver.Config, error) {
	return s.normalize(ctx, Config{})
}

func (s *configSchema) NormalizeSourceSpecific(
	ctx vcsdriver.ConfigNormalizeContext,
	def vcsdriver.Config,
) (vcsdriver.Config, error) {
	return s.normalize(ctx, def.(Config))
}

func (s *configSchema) normalize(
	ctx vcsdriver.ConfigNormalizeContext,
	cfg Config,
) (Config, error) {
	if s.SSHKey != nil {
		cfg.SSHKeyFile = s.SSHKey.File
		cfg.SSHKeyPassphrase = s.SSHKey.Passphrase

		if err := ctx.NormalizePath(&cfg.SSHKeyFile); err != nil {
			return Config{}, err
		}
	}

	if s.PreferHTTP != nil {
		cfg.PreferHTTP = *s.PreferHTTP
	}

	return cfg, nil
}
