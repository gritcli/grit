package sourcedriver

import (
	"context"

	"github.com/gritcli/grit/daemon/internal/logs"
)

// Cloner is an interface for cloning a specific remote repository.
//
// Cloners are obtained via the NewCloner() method on a Source.
type Cloner interface {
	// Clone makes a local clone of the remote repository in the given
	// directory.
	Clone(
		ctx context.Context,
		dir string,
		log logs.Log,
	) error
}
