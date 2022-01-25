package stubs

import (
	"errors"

	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2"
)

// VCSDriverConfigLoader is a test implementation of vcsdriver.ConfigLoader.
type VCSDriverConfigLoader struct {
	DefaultsFunc func(vcsdriver.ConfigContext) (vcsdriver.Config, error)
	MergeFunc    func(vcsdriver.ConfigContext, vcsdriver.Config, hcl.Body) (vcsdriver.Config, error)
}

// Defaults returns s.DefaultsFunc() if it is non-nil; otherwise, it returns an
// error.
func (s *VCSDriverConfigLoader) Defaults(
	ctx vcsdriver.ConfigContext,
) (vcsdriver.Config, error) {
	if s.DefaultsFunc != nil {
		return s.DefaultsFunc(ctx)
	}

	return nil, errors.New("<not implemented>")
}

// Merge returns s.MergeFunc() if it is non-nil; otherwise, it returns an error.
func (s *VCSDriverConfigLoader) Merge(
	ctx vcsdriver.ConfigContext,
	c vcsdriver.Config,
	b hcl.Body,
) (vcsdriver.Config, error) {
	if s.MergeFunc != nil {
		return s.MergeFunc(ctx, c, b)
	}

	return nil, errors.New("<not implemented>")
}

// VCSDriverConfigSchema is the HCL schema for VCSConfig.
type VCSDriverConfigSchema struct {
	ArbitraryAttribute string `hcl:"arbitrary_attribute,optional"`
	FilesystemPath     string `hcl:"filesystem_path,optional"`
}

// VCSDriverConfig is a test implementation of vcsdriver.Config.
type VCSDriverConfig struct {
	DescribeVCSConfigFunc func() string

	ArbitraryAttribute string
	FilesystemPath     string
}

// DescribeVCSConfig returns s.DescribeVCSConfigFunc() if it is non-nil;
// otherwise, it returns a new fixed value.
func (s *VCSDriverConfig) DescribeVCSConfig() string {
	if s.DescribeVCSConfigFunc != nil {
		return s.DescribeVCSConfigFunc()
	}

	return "<description>"
}
