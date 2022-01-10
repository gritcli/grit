package source

import (
	"context"
	"errors"
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
		return "", fmt.Errorf("unable to clone: unrecognized source (%s)", source)
	}

	bc, dir, err := src.Driver.NewBoundCloner(
		ctx,
		repoID,
		clientLogger,
	)
	if err != nil {
		return "", fmt.Errorf("unable to prepare for cloning: %w", err)
	}

	dir = filepath.Join(src.CloneDir, dir)

	if err := makeCloneDir(dir); err != nil {
		return "", fmt.Errorf("unable to create clone directory (%s): %w", dir, err)
	}

	if err := bc.Clone(ctx, dir); err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("unable to clone: %w", err)
	}

	return dir, nil
}

// makeCloneDir makes the given directory (and all of its parents) only if it
// does not already exist.
func makeCloneDir(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		return errors.New("file or directory already exists")
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(dir, 0700)
}
