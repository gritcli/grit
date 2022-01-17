package config

import (
	"fmt"

	"github.com/gritcli/grit/internal/common/api"
	"github.com/mitchellh/go-homedir"
)

// mergeDaemon merges s into the configuration.
func (r *resolver) mergeDaemon(s daemonSchema) error {
	if r.daemonFile != "" {
		return fmt.Errorf(
			"%s: the daemon configuration is already defined in %s",
			r.currentFile,
			r.daemonFile,
		)
	}

	d := Daemon(s)

	if err := normalizePath(r.currentFile, &d.Socket); err != nil {
		return err
	}

	r.daemonFile = r.currentFile
	r.output.Daemon = d

	return nil
}

// populateDaemonDefaults populates d with default values.
func (r *resolver) populateDaemonDefaults(d *Daemon) error {
	if d.Socket == "" {
		s, err := homedir.Expand(api.DefaultSocket)
		if err != nil {
			return fmt.Errorf(
				"unable to determine default daemon socket path: %w",
				err,
			)
		}

		d.Socket = s
	}

	return nil
}
