package source

import (
	"context"
)

// Source is an interface for a repository source.
type Source interface {
	// Name returns a short, human-readable identifier of the repository source.
	Name() string

	// Description returns a brief description of the repository source.
	Description() string

	// Init initializes the source.
	Init(ctx context.Context) error

	// Run runs any background processes required by the source until ctx is
	// canceled or a fatal error occurs.
	Run(ctx context.Context) error
}
