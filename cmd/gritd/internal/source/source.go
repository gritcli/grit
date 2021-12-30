package source

import (
	"context"
)

// Source is an interface for a repository source.
type Source interface {
	// Name returns a short, human-readable identifier of the repository source.
	Name() string

	// Description returns a brief description of the repository source.
	//
	// It may return limited information before the source has been initialized.
	Description() string

	// Init initializes the source.
	//
	// It is called before the daemon starts serving API requests.
	Init(ctx context.Context) error

	// Run performs any background tasks required by the source.
	//
	// It is called after the source is initialized.
	Run(ctx context.Context) error

	// Resolve resolves a repository name to a set of possible repositories.
	//
	// It does not perform partial-matching of the name, but may treat the name
	// as ambiguous by returning multiple repositories.
	//
	// l is a target for log messages to display to the user.
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
