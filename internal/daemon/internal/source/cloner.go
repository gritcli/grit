package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dogmatiq/dodeca/logging"
)

// BoundCloner is an interface for cloning repositories that is "bound" to a
// specific repository.
//
// A BoundCloner can be obtained by calling the Source.NewBoundCloner() method.
type BoundCloner interface {
	// Clone makes a local clone of the remote repository in the given
	// directory.
	Clone(ctx context.Context, dir string) error
}

// A Cloner clones repositories.
type Cloner struct {
	Sources List
	Logger  logging.Logger
}

// Clone clones a repository identified by source name and ID and returns the
// directory it was cloned into.
func (c *Cloner) Clone(
	ctx context.Context,
	source, repoID string,
	clientLogger logging.Logger,
) (string, error) {
	src, ok := c.Sources.ByName(source)
	if !ok {
		return "", fmt.Errorf("unrecognized source (%s)", source)
	}

	bc, dir, err := src.Driver.NewBoundCloner(
		ctx,
		repoID,
		clientLogger,
	)
	if err != nil {
		return "", err
	}

	dir = filepath.Join(src.CloneDir, dir)

	if _, err := os.Stat(dir); err == nil {
		return "", fmt.Errorf("clone directory (%s) already exists", dir)
	} else if !os.IsNotExist(err) {
		return "", err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	if err := bc.Clone(ctx, dir); err != nil {
		os.RemoveAll(dir) //nolint:errcheck
		return "", err
	}

	return dir, nil
}
