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
	Clone(
		ctx context.Context,
		dir string,
		logger logging.Logger,
	) error
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
) (_ string, err error) {
	logger := logging.Tee(
		clientLogger,
		logging.Prefix(
			c.Logger,
			"source[%s]: clone %s: ",
			source,
			repoID,
		),
	)

	defer func() {
		if err != nil {
			logger.LogString(err.Error())
		}
	}()

	src, ok := c.Sources.ByName(source)
	if !ok {
		return "", fmt.Errorf("unable to clone: unrecognized source (%s)", source)
	}

	bc, dir, err := src.Driver.NewBoundCloner(ctx, repoID, logger)
	if err != nil {
		return "", fmt.Errorf("unable to prepare for cloning: %w", err)
	}

	dir = filepath.Join(src.CloneDir, dir)

	if err := makeCloneDir(dir); err != nil {
		return "", fmt.Errorf("unable to create clone directory: %w", err)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	if err := bc.Clone(ctx, dir, logger); err != nil {
		return "", fmt.Errorf("unable to clone: %w", err)
	}

	return dir, nil
}

// makeCloneDir makes the given directory (and all of its parents) only if it
// does not already exist.
func makeCloneDir(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("mkdir %s: file exists", dir)
	} else if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0700)
	} else {
		return err
	}
}
