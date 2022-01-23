package stubs

import (
	"context"
	"errors"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
)

// SourceDriverConfigSchema is a test implementation of
// sourcedriver.ConfigSchema.
type SourceDriverConfigSchema struct {
	NormalizeFunc func(sourcedriver.ConfigNormalizeContext, *SourceDriverConfigSchema) (sourcedriver.Config, error)

	// These attributes must be defined in _this_ struct in order to use it as
	// the HCL schema.

	ArbitraryAttribute string `hcl:"arbitrary_attribute,optional"`
	FilesystemPath     string `hcl:"filesystem_path,optional"`
}

// Normalize returns s.NormalizeFunc() if it is non-nil, otherwise returns a
// new SourceDriverConfig stub.
func (s *SourceDriverConfigSchema) Normalize(
	nc sourcedriver.ConfigNormalizeContext,
) (sourcedriver.Config, error) {
	if s.NormalizeFunc != nil {
		return s.NormalizeFunc(nc, s)
	}

	return &SourceDriverConfig{}, nil
}

// SourceDriverConfig is a test implementation of the sourcedriver.Config
// interface.
type SourceDriverConfig struct {
	NewDriverFunc            func() sourcedriver.Driver
	DescribeSourceConfigFunc func() string

	ArbitraryAttribute string
	FilesystemPath     string
	VCSs               map[string]vcsdriver.Config
}

// NewDriver returns s.NewDriverFunc() if it is non-nil; otherwise, it returns a
// new SourceDriver stub.
func (s *SourceDriverConfig) NewDriver() sourcedriver.Driver {
	if s.NewDriverFunc != nil {
		return s.NewDriverFunc()
	}

	return &SourceDriver{}
}

// DescribeSourceConfig returns s.DescribeSourceConfigFunc() if it is non-nil;
// otherwise, it returns a new fixed value.
func (s *SourceDriverConfig) DescribeSourceConfig() string {
	if s.DescribeSourceConfigFunc != nil {
		return s.DescribeSourceConfigFunc()
	}

	return "<description>"
}

// SourceDriver is a test implementation of the sourcedriver.Driver interface.
type SourceDriver struct {
	InitFunc      func(context.Context, logging.Logger) error
	RunFunc       func(context.Context, logging.Logger) error
	StatusFunc    func(context.Context) (string, error)
	ResolveFunc   func(context.Context, string, logging.Logger) ([]sourcedriver.RemoteRepo, error)
	NewClonerFunc func(context.Context, string, logging.Logger) (sourcedriver.Cloner, string, error)
}

// Init returns s.InitFunc() if it is non-nil; otherwise, it returns nil.
func (s *SourceDriver) Init(ctx context.Context, logger logging.Logger) error {
	if s.InitFunc != nil {
		return s.InitFunc(ctx, logger)
	}

	return nil
}

// Run returns s.RunFunc() if it is non-nil; otherwise, it returns nil.
func (s *SourceDriver) Run(ctx context.Context, logger logging.Logger) error {
	if s.RunFunc != nil {
		return s.RunFunc(ctx, logger)
	}

	return nil
}

// Status returns s.StatusFunc() if it is non-nil; otherwise, it returns a fixed
// value.
func (s *SourceDriver) Status(ctx context.Context) (string, error) {
	if s.StatusFunc != nil {
		return s.StatusFunc(ctx)
	}

	return "<status>", nil
}

// Resolve returns s.ResolveFunc() if it is non-nil; otherwise, it returns a
// (nil, nil).
func (s *SourceDriver) Resolve(
	ctx context.Context,
	query string,
	logger logging.Logger,
) ([]sourcedriver.RemoteRepo, error) {
	if s.ResolveFunc != nil {
		return s.ResolveFunc(ctx, query, logger)
	}

	return nil, nil
}

// NewCloner returns s.NewClonerFunc() if it is non-nil; otherwise, it returns
// an error.
func (s *SourceDriver) NewCloner(
	ctx context.Context,
	id string,
	logger logging.Logger,
) (c sourcedriver.Cloner, dir string, err error) {
	if s.NewClonerFunc != nil {
		return s.NewClonerFunc(ctx, id, logger)
	}

	return nil, "", errors.New("<not implemented>")
}

// SourceDriverCloner is a test implementation of the sourcedriver.Cloner
// interface.
type SourceDriverCloner struct {
	CloneFunc func(context.Context, string, logging.Logger) error
}

// Clone returns s.CloneFunc() if it is non-nil; otherwise, it returns an error.
func (s *SourceDriverCloner) Clone(
	ctx context.Context,
	dir string,
	logger logging.Logger,
) error {
	if s.CloneFunc != nil {
		return s.CloneFunc(ctx, dir, logger)
	}

	return errors.New("<not implemented>")
}