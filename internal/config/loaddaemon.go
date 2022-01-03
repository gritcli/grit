package config

import "fmt"

// mergeDaemonBlock merges b into the configuration.
func (l *loader) mergeDaemonBlock(filename string, b daemonBlock) error {
	if l.daemonBlockFile != "" {
		return fmt.Errorf("the daemon configuration has already been defined in %s", l.daemonBlockFile)
	}

	d := Daemon(b)

	if d.Socket == "" {
		d.Socket = DefaultConfig.Daemon.Socket
	}

	if err := normalizePath(filename, &d.Socket); err != nil {
		return fmt.Errorf("the daemon socket path can not be resolved: %w", err)
	}

	l.daemonBlockFile = filename
	l.config.Daemon = d

	return nil
}
