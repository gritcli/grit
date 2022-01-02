package source

import (
	"context"
	"io"
)

// Source is an interface for a "repository source" that makes repositories
// available to Grit.
type Source interface {
	// Name returns a short, human-readable identifier of the repository source.
	Name() string

	// Description returns a brief description of the repository source.
	//
	// It may be called at any time, including before the source has been
	// initialized. It may return limited information until the source has been
	// initialized.
	Description() string

	// Init initializes the source.
	//
	// It is called before the daemon starts serving API requests.
	Init(ctx context.Context) error

	// Run performs any background tasks required by the source.
	//
	// It is called after the source is initialized and should run until ctx is
	// canceled or there is nothing left to do. The context is canceled when the
	// daemon shuts down.
	Run(ctx context.Context) error

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
	// query is invalid; an invalid query may be valid for other sources.
	//
	// out is a target for any output that should be displayed to the user
	Resolve(ctx context.Context, query string, out io.Writer) ([]Repo, error)
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
