package config

import (
	"fmt"

	"github.com/gritcli/grit/internal/common/api"
	homedir "github.com/mitchellh/go-homedir"
)

// mergeDaemonBlock merges a "daemon" block into the configuration.
func mergeDaemonBlock(
	cfg *unresolvedConfig,
	filename string,
	daemon daemonSchema,
) error {
	if cfg.Daemon.File != "" {
		return fmt.Errorf(
			"%s: a 'daemon' block is already defined in %s",
			filename,
			cfg.Daemon.File,
		)
	}

	cfg.Daemon.File = filename
	cfg.Daemon.Block = daemon

	return nil
}

// normalizeDaemonBlock normalizes cfg.Daemon.Block and populates it with
// default values.
func normalizeDaemonBlock(cfg *unresolvedConfig) error {
	if cfg.Daemon.Block.Socket != "" {
		return normalizePath(cfg.Daemon.File, &cfg.Daemon.Block.Socket)
	}

	dir, err := homedir.Expand(api.DefaultSocket)
	if err != nil {
		return err
	}

	cfg.Daemon.Block.Socket = dir

	return nil
}
