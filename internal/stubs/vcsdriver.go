package stubs

import (
	"github.com/gritcli/grit/driver/vcsdriver"
)

// VCSDriverConfigSchema is a test implementation of vcsdriver.ConfigSchema.
type VCSDriverConfigSchema struct {
	NormalizeGlobalsFunc        func(vcsdriver.ConfigNormalizeContext, *VCSDriverConfigSchema) (vcsdriver.Config, error)
	NormalizeSourceSpecificFunc func(vcsdriver.ConfigNormalizeContext, vcsdriver.Config, *VCSDriverConfigSchema) (vcsdriver.Config, error)

	// These attributes must be defined in _this_ struct in order to use it as
	// the HCL schema.

	ArbitraryAttribute string `hcl:"arbitrary_attribute,optional"`
	FilesystemPath     string `hcl:"filesystem_path,optional"`
}

// NormalizeGlobals returns s.NormalizeGlobalsFunc() if it is non-nil, otherwise
// returns a new VCSDriverConfig stub.
func (s *VCSDriverConfigSchema) NormalizeGlobals(
	nc vcsdriver.ConfigNormalizeContext,
) (vcsdriver.Config, error) {
	if s.NormalizeGlobalsFunc != nil {
		return s.NormalizeGlobalsFunc(nc, s)
	}

	return &VCSDriverConfig{}, nil
}

// NormalizeSourceSpecific returns s.NormalizeSourceSpecificFunc() if it is
// non-nil, otherwise returns g.
func (s *VCSDriverConfigSchema) NormalizeSourceSpecific(
	nc vcsdriver.ConfigNormalizeContext,
	g vcsdriver.Config,
) (vcsdriver.Config, error) {
	if s.NormalizeSourceSpecificFunc != nil {
		return s.NormalizeSourceSpecificFunc(nc, g, s)
	}

	return g, nil
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
