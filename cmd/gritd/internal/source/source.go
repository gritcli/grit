package source

import (
	"context"
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

	// Resolve resolves a repository name to a set of possible repositories.
	//
	// The exact name-matching logic is implementation defined. Implementations
	// should be as generous as possible in what they accept but should avoid
	// returning repositories with names that only partially match the input.
	// Multiple repositories may be returned to indicate that the name is
	// ambiguous.
	//
	// The name is typically captured directly from user input and has not been
	// sanitized. The implementation must not return an error if the name is
	// invalid; an invalid name may be valid for other sources.
	Resolve(ctx context.Context, name string) ([]Repo, error)
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
