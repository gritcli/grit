package source

import "context"

// BoundCloner is an interface for cloning repositories that is "bound" to a
// specific repository.
//
// A BoundCloner can be obtained by calling the Source.NewBoundCloner() method.
type BoundCloner interface {
	// Clone makes a local clone of the remote repository in the given
	// directory.
	Clone(ctx context.Context, dir string) error
}
