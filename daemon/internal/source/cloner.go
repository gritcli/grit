package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
)

// LocalRepo represents a local clone of a remote repository.
type LocalRepo struct {
	sourcedriver.RemoteRepo
	Source Source

	// AbsoluteCloneDir is the absolute path to the directory containing the
	// local clone.
	AbsoluteCloneDir string
}

// A Cloner clones repositories.
type Cloner struct {
	Sources List
	Log     logs.Log
}

// Clone clones a repository identified by source name and ID and returns the
// directory it was cloned into.
func (c *Cloner) Clone(
	ctx context.Context,
	source, repoID string,
	clientLog logs.Log,
) (_ LocalRepo, err error) {
	log := logs.Tee(
		clientLog,
		newLog(source, c.Log).
			WithPrefix("clone %s: ", repoID),
	)

	defer func() {
		if err != nil {
			log.Write("%s", err.Error())
		}
	}()

	src, ok := c.Sources.ByName(source)
	if !ok {
		return LocalRepo{}, fmt.Errorf("unable to clone: unrecognized source (%s)", source)
	}

	cloner, repo, err := src.Driver.Cloner(ctx, repoID, log)
	if err != nil {
		return LocalRepo{}, fmt.Errorf("unable to prepare for cloning: %w", err)
	}

	dir := filepath.Join(src.BaseCloneDir, repo.RelativeCloneDir)

	if err := makeCloneDir(dir); err != nil {
		return LocalRepo{}, fmt.Errorf("unable to create clone directory: %w", err)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	if err := cloner.Clone(ctx, dir, log); err != nil {
		return LocalRepo{}, fmt.Errorf("unable to clone: %w", err)
	}

	return LocalRepo{
		repo,
		src,
		filepath.Join(
			src.BaseCloneDir,
			repo.RelativeCloneDir,
		),
	}, nil
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
