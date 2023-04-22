package stubs

import (
	"context"
	"errors"
	"net/http"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
	"github.com/gritcli/grit/daemon/internal/logs"
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
	InitFunc      func(context.Context, logs.Log) error
	RunFunc       func(context.Context, logs.Log) error
	StatusFunc    func(context.Context, logs.Log) (string, error)
	SignInFunc    func(context.Context, logs.Log) error
	SignOutFunc   func(context.Context, logs.Log) error
	ResolveFunc   func(context.Context, string, logs.Log) ([]sourcedriver.RemoteRepo, error)
	ClonerFunc    func(context.Context, string, logs.Log) (sourcedriver.Cloner, sourcedriver.RemoteRepo, error)
	SuggestFunc   func(string, logs.Log) map[string][]sourcedriver.RemoteRepo
	ServeHTTPFunc http.HandlerFunc
}

// Init returns s.InitFunc() if it is non-nil; otherwise, it returns nil.
func (s *Source) Init(ctx context.Context, log logs.Log) error {
	if s.InitFunc != nil {
		return s.InitFunc(ctx, log)
	}

	return nil
}

// Run returns s.RunFunc() if it is non-nil; otherwise, it returns nil.
func (s *Source) Run(ctx context.Context, log logs.Log) error {
	if s.RunFunc != nil {
		return s.RunFunc(ctx, log)
	}

	return nil
}

// Status returns s.StatusFunc() if it is non-nil; otherwise, it returns a fixed
// value.
func (s *Source) Status(ctx context.Context, log logs.Log) (string, error) {
	if s.StatusFunc != nil {
		return s.StatusFunc(ctx, log)
	}

	return "<status>", nil
}

// SignIn returns s.SignInFunc() if it is non-nil; otherwise, it
// returns nil.
func (s *Source) SignIn(ctx context.Context, log logs.Log) error {
	if s.SignInFunc != nil {
		return s.SignInFunc(ctx, log)
	}

	return nil
}

// SignOut returns s.SignOutFunc() if it is non-nil; otherwise, it returns nil.
func (s *Source) SignOut(ctx context.Context, log logs.Log) error {
	if s.SignOutFunc != nil {
		return s.SignOutFunc(ctx, log)
	}
	return nil
}

// Resolve returns s.ResolveFunc() if it is non-nil; otherwise, it returns a
// (nil, nil).
func (s *Source) Resolve(
	ctx context.Context,
	query string,
	log logs.Log,
) ([]sourcedriver.RemoteRepo, error) {
	if s.ResolveFunc != nil {
		return s.ResolveFunc(ctx, query, log)
	}

	return nil, nil
}

// Cloner returns s.ClonerFunc() if it is non-nil; otherwise, it returns an
// error.
func (s *Source) Cloner(
	ctx context.Context,
	id string,
	log logs.Log,
) (sourcedriver.Cloner, sourcedriver.RemoteRepo, error) {
	if s.ClonerFunc != nil {
		return s.ClonerFunc(ctx, id, log)
	}

	return nil, sourcedriver.RemoteRepo{}, errors.New("<not implemented>")
}

// Suggest returns s.SuggestFunc() if it is non-nil; otherwise, it returns nil.
func (s *Source) Suggest(word string, log logs.Log) map[string][]sourcedriver.RemoteRepo {
	if s.SuggestFunc != nil {
		return s.SuggestFunc(word, log)
	}

	return nil
}

// ServeHTTP calls s.ServeHTTPFunc() if it is non-nil.
func (s *Source) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.ServeHTTPFunc != nil {
		s.ServeHTTPFunc(w, r)
	}
}

// SourceCloner is a test implementation of the sourcedriver.SourceCloner
// interface.
type SourceCloner struct {
	CloneFunc func(context.Context, string, logs.Log) error
}

// Clone returns s.CloneFunc() if it is non-nil; otherwise, it returns nil.
func (s *SourceCloner) Clone(
	ctx context.Context,
	dir string,
	log logs.Log,
) error {
	if s.CloneFunc != nil {
		return s.CloneFunc(ctx, dir, log)
	}

	return nil
}
