package source

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/daemon/internal/scm"
)

// Source is a repository source.
type Source struct {
	// Name is the unique name for the repository source.
	Name string

	// Description is a human-readable description of the source.
	Description string

	// CloneDir is the directory containing repositories cloned from this
	// source.
	CloneDir string

	// Driver is the driver used to perform repository operations.
	Driver Driver
}

// A Driver performs implementation-specific repository operations for a
// repository source.
type Driver interface {
	// Init initializes the driver.
	//
	// It is called before the daemon starts serving API requests. If an error
	// is returned, the daemon fails to start.
	Init(ctx context.Context) error

	// Run performs any ongoing behavior required by the driver.
	//
	// It is called in its own goroutine after the driver is initialized, It
	// should run until ctx is canceled or there is nothing left to do. The
	// context is canceled when the daemon shuts down.
	//
	// If it returns an error before ctx is canceled, the daemon is stopped.
	Run(ctx context.Context) error

	// Status returns a brief description of the status of the driver.
	//
	// This may include information about the driver's ability to communicate
	// with the remote server, the authenticated user, etc.
	Status(ctx context.Context) (string, error)

	// Resolve resolves a repository name, URL, or other identifier to a set of
	// possible repositories.
	//
	// The details of the resolution logic is implementation defined.
	// Implementations should be as generous as possible in what they accept but
	// should avoid returning repositories that only partially match the input.
	// Multiple repositories may be returned to indicate that the query is
	// ambiguous.
	//
	// The query string is typically captured directly from user input and has
	// not been sanitized. The implementation must not return an error if the
	// query is invalid; instead return an empty slice.
	//
	// clientLog is a target for any log output that should be sent to the
	// client and displayed to the user.
	Resolve(ctx context.Context, query string, clientLog logging.Logger) ([]Repo, error)

	// NewCloner returns an scm.Cloner that clones the specified repository.
	//
	// id is the repository ID, as discovered by a prior call to Resolve().
	//
	// clientLog is a target for any log output that should be sent to the
	// client and displayed to the user while cloning.
	//
	// dir is the sub-directory that the clone should be placed into, relative
	// to the source's configured clone directory. Typically this should be some
	// form of the repository's name, sanitized for use as a directory name.
	NewCloner(
		ctx context.Context,
		id string,
		clientLog logging.Logger,
	) (c scm.Cloner, dir string, err error)
}

// Repo is a reference to a remote repository provided by a source.
type Repo struct {
	// ID is a unique (within the source) identifier for the repository.
	ID string

	// Name is the name of the repository.
	Name string

	// Description is a human-readable description of the repository.
	Description string

	// WebURL is the URL of the repository's web page, if available.
	WebURL string
}
