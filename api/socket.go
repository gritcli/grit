package api

import "path/filepath"

// DefaultSocket is the default path the Unix socket used for
// communication between the CLI and the daemon.
//
// Even though this is a "Unix socket", the AF_UNIX address family is
// supported on Windows 10+.
var DefaultSocket = filepath.Join("~", "grit", "daemon.sock")
