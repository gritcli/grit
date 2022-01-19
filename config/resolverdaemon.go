package config

import (
	"fmt"

	"github.com/gritcli/grit/internal/common/api"
)

// mergeDaemon merges s into the configuration.
func (r *resolver) mergeDaemon(file string, s daemonSchema) error {
	if r.daemonFile != "" {
		return fmt.Errorf(
			"%s: the daemon configuration is already defined in %s",
			file,
			r.daemonFile,
		)
	}

	cfg := Daemon(s)

	if err := r.normalizePath(&cfg.Socket); err != nil {
		return err // TODO: explain error!
	}

	r.daemonFile = file
	r.daemon = cfg

	return nil
}

// populateDaemonDefaults populates r.daemon with default values.
// TODO: can this be moved into the mergeDaemon() function?
func (r *resolver) populateDaemonDefaults() error {
	if r.daemon.Socket == "" {
		r.daemon.Socket = api.DefaultSocket

		if err := r.normalizePath(&r.daemon.Socket); err != nil {
			return err // TODO: explain error!
		}
	}

	return nil
}
