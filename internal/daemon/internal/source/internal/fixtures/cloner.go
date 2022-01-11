package fixtures

import (
	"context"
	"errors"

	"github.com/dogmatiq/dodeca/logging"
)

// BoundClonerStub is a test implementation of the source.BoundCloner interface.
type BoundClonerStub struct {
	CloneFunc func(context.Context, string, logging.Logger) error
}

// Clone calls s.CloneFunc(ctx, dir, logger) if s.CloneFunc is non-nil;
// otherwise it returns an error.
func (s *BoundClonerStub) Clone(
	ctx context.Context,
	dir string,
	logger logging.Logger,
) error {
	if s.CloneFunc != nil {
		return s.CloneFunc(ctx, dir, logger)
	}

	return errors.New("<not implemented>")
}
