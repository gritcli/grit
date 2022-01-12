package config

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// Daemon holds the configuration for the Grit daemon.
type Daemon struct {
	// Socket is the path of the Unix socket used for communication between
	// the Grit CLI and the Grit daemon (via gRPC).
	Socket string
}

// daemonBlock is the HCL schema for a "daemon" block.
type daemonBlock struct {
	Socket string `hcl:"socket,optional"`
}

// mergeDaemonBlock merges b into cfg.
func mergeDaemonBlock(cfg *unresolvedConfig, filename string, b daemonBlock) error {
	if cfg.Daemon.File != "" {
		return fmt.Errorf(
			"%s: a 'daemon' block is already defined in %s",
			filename,
			cfg.Daemon.File,
		)
	}

	cfg.Daemon.File = filename
	cfg.Daemon.Block = b

	return nil
}

// normalizeDaemonBlock normalizes cfg.Daemon.Block and populates it with
// default values.
func normalizeDaemonBlock(cfg *unresolvedConfig) error {
	if cfg.Daemon.Block.Socket != "" {
		return normalizePath(cfg.Daemon.File, &cfg.Daemon.Block.Socket)
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("unable to determine the current user's home directory")
	}

	cfg.Daemon.Block.Socket = filepath.Join(homeDir, "grit", "daemon.sock")

	return nil
}

// assembleDaemonBlock converts b into its configuration representation.
func assembleDaemonBlock(b daemonBlock) Daemon {
	return Daemon(b)
}
