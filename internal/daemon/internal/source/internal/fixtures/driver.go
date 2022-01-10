package fixtures

import (
	"context"
	"errors"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/daemon/internal/source"
)

// DriverStub is a test implementation of the source.Driver interface.
type DriverStub struct {
	InitFunc           func(context.Context) error
	RunFunc            func(context.Context) error
	StatusFunc         func(context.Context) (string, error)
	ResolveFunc        func(context.Context, string, logging.Logger) ([]source.Repo, error)
	NewBoundClonerFunc func(context.Context, string, logging.Logger) (source.BoundCloner, string, error)
}

// Init calls s.InitFunc(ctx) if s.InitFunc is non-nil; otherwise it returns
// nil.
func (s *DriverStub) Init(ctx context.Context) error {
	if s.InitFunc != nil {
		return s.InitFunc(ctx)
	}

	return nil
}

// Run calls s.RunFunc(ctx) if s.RunFunc is non-nil; otherwise it returns nil.
func (s *DriverStub) Run(ctx context.Context) error {
	if s.RunFunc != nil {
		return s.RunFunc(ctx)
	}

	return nil
}

// Status calls s.StatusFunc(ctx) if s.StatusFunc is non-nil; otherwise it
// returns a default status message.
func (s *DriverStub) Status(ctx context.Context) (string, error) {
	if s.StatusFunc != nil {
		return s.StatusFunc(ctx)
	}

	return "<status>", nil
}

// Resolve calls s.ResolveFunc(ctx, query, clientLog) if s.ResolveFunc is
// non-nil; otherwise it returns (nil, nil).
func (s *DriverStub) Resolve(
	ctx context.Context,
	query string,
	clientLog logging.Logger,
) ([]source.Repo, error) {
	if s.ResolveFunc != nil {
		return s.ResolveFunc(ctx, query, clientLog)
	}

	return nil, nil
}

// NewBoundCloner calls s.NewBoundClonerFunc(ctx, id, clientLog) if
// s.NewBoundClonerFunc is non-nil; otherwise it returns an error.
func (s *DriverStub) NewBoundCloner(
	ctx context.Context,
	id string,
	clientLog logging.Logger,
) (c source.BoundCloner, dir string, err error) {
	if s.NewBoundClonerFunc != nil {
		return s.NewBoundClonerFunc(ctx, id, clientLog)
	}

	return nil, "", errors.New("<not implemented>")
}
