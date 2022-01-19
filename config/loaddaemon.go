package config

import (
	"fmt"
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
		return fmt.Errorf(
			"unable to resolve daemon socket path: %w (%s)",
			err,
			cfg.Socket,
		)
	}

	l.daemonFile = file
	l.daemon = cfg

	return nil
}

// populateDaemonDefaults populates l.daemon with default values.
func (l *loader) populateDaemonDefaults() error {
	if l.daemon.Socket == "" {
		l.daemon.Socket = DefaultDaemonSocket

		if err := l.normalizePath(&l.daemon.Socket); err != nil {
			return fmt.Errorf(
				"unable to resolve default daemon socket path: %w (%s)",
				err,
				l.daemon.Socket,
			)
		}
	}

	return nil
}
