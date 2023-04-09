package stubs

import (
	"errors"

	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2"
)

// VCSConfigLoader is a test implementation of vcsdriver.ConfigLoader.
type VCSConfigLoader struct {
	DefaultsFunc          func(vcsdriver.ConfigContext) (vcsdriver.Config, error)
	UnmarshalAndMergeFunc func(vcsdriver.ConfigContext, vcsdriver.Config, hcl.Body) (vcsdriver.Config, error)
}

// Defaults returns s.DefaultsFunc() if it is non-nil; otherwise, it returns an
// error.
func (s *VCSConfigLoader) Defaults(
	ctx vcsdriver.ConfigContext,
) (vcsdriver.Config, error) {
	if s.DefaultsFunc != nil {
		return s.DefaultsFunc(ctx)
	}

	return nil, errors.New("<not implemented>")
}

// UnmarshalAndMerge returns s.MergeFunc() if it is non-nil; otherwise, it
// returns an error.
func (s *VCSConfigLoader) UnmarshalAndMerge(
	ctx vcsdriver.ConfigContext,
	c vcsdriver.Config,
	b hcl.Body,
) (vcsdriver.Config, error) {
	if s.UnmarshalAndMergeFunc != nil {
		return s.UnmarshalAndMergeFunc(ctx, c, b)
	}

	return nil, errors.New("<not implemented>")
}

// VCSConfigSchema is the HCL schema for VCSConfig.
type VCSConfigSchema struct {
	ArbitraryAttribute string `hcl:"arbitrary_attribute,optional"`
	FilesystemPath     string `hcl:"filesystem_path,optional"`
}

// VCSConfig is a test implementation of vcsdriver.Config.
type VCSConfig struct {
	DescribeVCSConfigFunc func() string

	ArbitraryAttribute string
	FilesystemPath     string
}

// DescribeVCSConfig returns s.DescribeVCSConfigFunc() if it is non-nil;
// otherwise, it returns a new fixed value.
func (s *VCSConfig) DescribeVCSConfig() string {
	if s.DescribeVCSConfigFunc != nil {
		return s.DescribeVCSConfigFunc()
	}

	return "<description>"
}
