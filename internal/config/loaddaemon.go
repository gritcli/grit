package config

// resolve returns the configuration that is produced by this block.
func (b daemonBlock) resolve(filename string) (Daemon, error) {
	d := Daemon(b)

	if d.Socket == "" {
		d.Socket = DefaultConfig.Daemon.Socket
	}

	return d, normalizePath(filename, &d.Socket)
}
