package stubs

import (
	"context"
	"errors"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2"
)

// SourceConfigLoader is a test implementation of sourcedriver.ConfigLoader.
type SourceConfigLoader struct {
	UnmarshalFunc       func(sourcedriver.ConfigContext, hcl.Body) (sourcedriver.Config, error)
	ImplicitSourcesFunc func(sourcedriver.ConfigContext) ([]sourcedriver.ImplicitSource, error)
}

// Unmarshal returns s.UnmarshalFunc() if it is non-nil; otherwise, it returns
// an error.
func (s *SourceConfigLoader) Unmarshal(
	ctx sourcedriver.ConfigContext,
	b hcl.Body,
) (sourcedriver.Config, error) {
	if s.UnmarshalFunc != nil {
		return s.UnmarshalFunc(ctx, b)
	}

	return nil, errors.New("<not implemented>")
}

// ImplicitSources returns s.ImplicitSourcesFunc() if it is non-nil; otherwise,
// it returns (nil, nil).
func (s *SourceConfigLoader) ImplicitSources(
	ctx sourcedriver.ConfigContext,
) ([]sourcedriver.ImplicitSource, error) {
	if s.ImplicitSourcesFunc != nil {
		return s.ImplicitSourcesFunc(ctx)
	}

	return nil, nil
}

// SourceConfigSchema is the HCL schema for SourceConfig.
type SourceConfigSchema struct {
	ArbitraryAttribute string `hcl:"arbitrary_attribute,optional"`
	FilesystemPath     string `hcl:"filesystem_path,optional"`
}

// SourceConfig is a test implementation of the sourcedriver.Config interface.
type SourceConfig struct {
	NewSourceFunc            func() sourcedriver.Source
	DescribeSourceConfigFunc func() string

	ArbitraryAttribute string
	FilesystemPath     string
	VCSs               map[string]vcsdriver.Config
}

// NewSource returns s.NewSourceFunc() if it is non-nil; otherwise, it returns a
// new SourceDriver stub.
func (s *SourceConfig) NewSource() sourcedriver.Source {
	if s.NewSourceFunc != nil {
		return s.NewSourceFunc()
	}

	return &Source{}
}

// DescribeSourceConfig returns s.DescribeSourceConfigFunc() if it is non-nil;
// otherwise, it returns a new fixed value.
func (s *SourceConfig) DescribeSourceConfig() string {
	if s.DescribeSourceConfigFunc != nil {
		return s.DescribeSourceConfigFunc()
	}

	return "<description>"
}

// Source is a test implementation of the sourcedriver.Source interface.
type Source struct {
	InitFunc      func(context.Context, logging.Logger) error
	RunFunc       func(context.Context, logging.Logger) error
	StatusFunc    func(context.Context) (string, error)
	ResolveFunc   func(context.Context, string, logging.Logger) ([]sourcedriver.RemoteRepo, error)
	NewClonerFunc func(context.Context, string, logging.Logger) (sourcedriver.Cloner, sourcedriver.RemoteRepo, error)
	SuggestFunc   func(string) []sourcedriver.RemoteRepo
}

// Init returns s.InitFunc() if it is non-nil; otherwise, it returns nil.
func (s *Source) Init(ctx context.Context, logger logging.Logger) error {
	if s.InitFunc != nil {
		return s.InitFunc(ctx, logger)
	}

	return nil
}

// Run returns s.RunFunc() if it is non-nil; otherwise, it returns nil.
func (s *Source) Run(ctx context.Context, logger logging.Logger) error {
	if s.RunFunc != nil {
		return s.RunFunc(ctx, logger)
	}

	return nil
}

// Status returns s.StatusFunc() if it is non-nil; otherwise, it returns a fixed
// value.
func (s *Source) Status(ctx context.Context) (string, error) {
	if s.StatusFunc != nil {
		return s.StatusFunc(ctx)
	}

	return "<status>", nil
}

// Resolve returns s.ResolveFunc() if it is non-nil; otherwise, it returns a
// (nil, nil).
func (s *Source) Resolve(
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
func (s *Source) NewCloner(
	ctx context.Context,
	id string,
	logger logging.Logger,
) (sourcedriver.Cloner, sourcedriver.RemoteRepo, error) {
	if s.NewClonerFunc != nil {
		return s.NewClonerFunc(ctx, id, logger)
	}

	return nil, sourcedriver.RemoteRepo{}, errors.New("<not implemented>")
}

// Suggest returns s.SuggestFunc() if it is non-nil; otherwise, it returns nil.
func (s *Source) Suggest(word string) []sourcedriver.RemoteRepo {
	if s.SuggestFunc != nil {
		return s.SuggestFunc(word)
	}

	return nil
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
