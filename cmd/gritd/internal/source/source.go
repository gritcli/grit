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
}
