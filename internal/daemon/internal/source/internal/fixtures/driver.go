package fixtures

import (
	"context"
	"errors"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/plugin/driver"
)

// DriverConfigStub is a test implementation of the driver.Config interface.
type DriverConfigStub struct {
	NewDriverFunc func() driver.Driver
	StringFunc    func() string
}

// NewDriver returns s.NewDriverFunc() if s.NewDriverFunc is non-nil; otherwise
// it returns a new DriverStub.
func (s *DriverConfigStub) NewDriver() driver.Driver {
	if s.NewDriverFunc != nil {
		return s.NewDriverFunc()
	}

	return &DriverStub{}
}

// s.DriverConfigStub returns s.StringFunc() if s.StringFunc is non-nil;
// otherwise it returns a fixed string.
func (s *DriverConfigStub) String() string {
	if s.StringFunc != nil {
		return s.StringFunc()
	}

	return "<driver config stub>"
}

// DriverStub is a test implementation of the driver.Driver interface.
type DriverStub struct {
	InitFunc      func(context.Context, logging.Logger) error
	RunFunc       func(context.Context, logging.Logger) error
	StatusFunc    func(context.Context) (string, error)
	ResolveFunc   func(context.Context, string, logging.Logger) ([]driver.RemoteRepo, error)
	NewClonerFunc func(context.Context, string, logging.Logger) (driver.Cloner, string, error)
}

// Init calls s.InitFunc(ctx, logger) if s.InitFunc is non-nil; otherwise it
// returns nil.
func (s *DriverStub) Init(ctx context.Context, logger logging.Logger) error {
	if s.InitFunc != nil {
		return s.InitFunc(ctx, logger)
	}

	return nil
}

// Run calls s.RunFunc(ctx, logger) if s.RunFunc is non-nil; otherwise it
// returns nil.
func (s *DriverStub) Run(ctx context.Context, logger logging.Logger) error {
	if s.RunFunc != nil {
		return s.RunFunc(ctx, logger)
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

// Resolve calls s.ResolveFunc(ctx, query, logger) if s.ResolveFunc is non-nil;
// otherwise it returns (nil, nil).
func (s *DriverStub) Resolve(
	ctx context.Context,
	query string,
	logger logging.Logger,
) ([]driver.RemoteRepo, error) {
	if s.ResolveFunc != nil {
		return s.ResolveFunc(ctx, query, logger)
	}

	return nil, nil
}

// NewCloner calls s.NewClonerFunc(ctx, id, logger) if
// s.NewClonerFunc is non-nil; otherwise it returns an error.
func (s *DriverStub) NewCloner(
	ctx context.Context,
	id string,
	logger logging.Logger,
) (c driver.Cloner, dir string, err error) {
	if s.NewClonerFunc != nil {
		return s.NewClonerFunc(ctx, id, logger)
	}

	return nil, "", errors.New("<not implemented>")
}
