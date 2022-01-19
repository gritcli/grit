package config

import (
	"fmt"

	"github.com/gritcli/grit/internal/common/api"
)

// mergeDaemon merges s into the configuration.
func (l *loader) mergeDaemon(file string, s daemonSchema) error {
	if l.daemonFile != "" {
		return fmt.Errorf(
			"the daemon configuration is already defined in %s",
			l.daemonFile,
		)
	}

	cfg := Daemon(s)

	if err := l.normalizePath(&cfg.Socket); err != nil {
		return err // TODO: explain error!
	}

	l.daemonFile = file
	l.daemon = cfg

	return nil
}

// populateDaemonDefaults populates l.daemon with default values.
func (l *loader) populateDaemonDefaults() error {
	if l.daemon.Socket == "" {
		l.daemon.Socket = api.DefaultSocket

		if err := l.normalizePath(&l.daemon.Socket); err != nil {
			return err // TODO: explain error!
		}
	}

	return nil
}
