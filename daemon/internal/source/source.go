package source

import (
	"net/url"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
)

// Source is a repository source.
type Source struct {
	// Name is the unique name for the repository source.
	Name string

	// Description is a human-readable description of the source.
	Description string

	// BaseCloneDir is the directory containing repositories cloned from this
	// source.
	BaseCloneDir string

	// BaseURL is the base URL for that the daemon's HTTP server route's to the
	// source's HTTP handler implementation.
	BaseURL *url.URL

	// Driver is the source implementation provided by the driver, used to
	// perform repository operations for this source.
	Driver sourcedriver.Source
}

// Log returns the logger to use for messages about this source.
func (s Source) Log(log logs.Log) logs.Log {
	return newLog(s.Name, log)
}

func newLog(name string, log logs.Log) logs.Log {
	return log.WithPrefix("source/%s: ", name)
}
