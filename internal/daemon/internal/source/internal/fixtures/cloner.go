package fixtures

import (
	"context"
	"errors"
)

// BoundClonerStub is a test implementation of the source.BoundCloner interface.
type BoundClonerStub struct {
	CloneFunc func(context.Context, string) error
}

// Clone calls s.CloneFunc(ctx, dir) if s.CloneFunc is non-nil; otherwise it
// returns an error.
func (s *BoundClonerStub) Clone(ctx context.Context, dir string) error {
	if s.CloneFunc != nil {
		return s.CloneFunc(ctx, dir)
	}

	return errors.New("<not implemented>")
}
