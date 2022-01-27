package sourcedriver

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
)

// Source is an interface for a source provided by this driver.
type Source interface {
	// Init initializes the source.
	//
	// It is called before the daemon starts serving API requests. If an error
	// is returned, the daemon is stopped.
	Init(ctx context.Context, logger logging.Logger) error

	// Run performs any background processing required by the source.
	//
	// It is called in its own goroutine after the source is initialized. It
	// should run until ctx is canceled or there is nothing left to do. The
	// context is canceled when the daemon shuts down.
	//
	// If it returns an error before ctx is canceled, the daemon is stopped.
	Run(ctx context.Context, logger logging.Logger) error

	// Status returns a brief description of the current state of the source.
	//
	// This may include information about connectivity with a remote server,
	// authenticated details, etc.
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
	Resolve(
		ctx context.Context,
		query string,
		logger logging.Logger,
	) ([]RemoteRepo, error)

	// NewCloner returns a cloner that clones the repository with the given ID,
	// and information about the repository being cloned.
	//
	// id is the repository ID, as discovered by a prior call to Resolve().
	NewCloner(
		ctx context.Context,
		id string,
		logger logging.Logger,
	) (Cloner, RemoteRepo, error)

	// Suggest returns a set of repositories that have names beginning with the
	// given word (which may be empty).
	Suggest(word string) []RemoteRepo
}
