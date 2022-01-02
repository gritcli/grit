package github

import (
	"context"
	"errors"

	"github.com/dogmatiq/dodeca/logging"
)

// Clone makes a repository available at the specified directory.
func (s *impl) Clone(
	ctx context.Context,
	repoID, dir string,
	out logging.Logger,
) error {
	return errors.New("not implemented")
}
