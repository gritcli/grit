package config_test

import (
	"errors"

	. "github.com/gritcli/grit/daemon/internal/config"
	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
	"github.com/gritcli/grit/daemon/internal/stubs"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("func Load() (VCS configuration)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"sources inherit the global VCS configuration",
			[]string{
				`vcs "test_vcs_driver" {
					arbitrary_attribute = "<explicit>"
				}

				source "test_source" "test_source_driver" {}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<default>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default> + <explicit>",
						},
					},
				},
			}),
		),
		Entry(
			"sources can specify a VCS configuration",
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "test_vcs_driver" {
						arbitrary_attribute = "<explicit>"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<default>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default> + <explicit>",
						},
					},
				},
			}),
		),
		Entry(
			"sources can override the global VCS configuration",
			[]string{
				`vcs "test_vcs_driver" {
					arbitrary_attribute = "<explicit global>"
				}

				source "test_source" "test_source_driver" {
					vcs "test_vcs_driver" {
						arbitrary_attribute = "<override>"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<default>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default> + <explicit global> + <override>",
						},
					},
				},
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the global configuration",
		testLoadFailure,
		Entry(
			`empty VCS driver name`,
			[]string{
				`vcs "" {}`,
			},
			`<dir>/config-0.hcl: global VCS configuration with empty driver name`,
		),
		Entry(
			`duplicate global VCS configuration`,
			[]string{
				`vcs "test_vcs_driver" {}`,
				`vcs "test_vcs_driver" {}`,
			},
			`<dir>/config-1.hcl: global configuration for the 'test_vcs_driver' version control system is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`unrecognized VCS driver name`,
			[]string{
				`vcs "<unrecognized>" {}`,
			},
			`<dir>/config-0.hcl: the '<unrecognized>' version control system is not unrecognized, the supported VCS drivers are: 'test_vcs_driver'`,
		),
		Entry(
			`VCS defaults with a well-structured, but invalid body`,
			[]string{
				`vcs "test_vcs_driver" {
					unrecognized = true
				}`,
			},
			`<dir>/config-0.hcl:2,6-18: Unsupported argument; An argument named "unrecognized" is not expected here.`,
		),
		Entry(
			`error normalizing driver configuration`,
			[]string{
				`vcs "test_vcs_driver" {
					filesystem_path = "~someuser/path/to/nowhere"
				}`,
			},
			`<dir>/config-0.hcl: the global configuration for the 'test_vcs_driver' version control system cannot be loaded: cannot expand user-specific home dir`,
		),
		Entry(
			`error normalizing default driver configuration`,
			[]string{},
			`unable to produce default global configuration for the 'test_vcs_driver_with_broken_default' version control system: <error>`,
			func(reg *DriverRegistry) {
				reg.RegisterVCSDriver(
					"test_vcs_driver_with_broken_default",
					vcsdriver.Registration{
						Name: testVCSDriverName,
						ConfigLoader: &stubs.VCSConfigLoader{
							DefaultsFunc: func(
								vcsdriver.ConfigContext,
							) (vcsdriver.Config, error) {
								return nil, errors.New("<error>")
							},
						},
					},
				)
			},
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with a source-specific configuration",
		testLoadFailure,
		Entry(
			`empty VCS driver name`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "" {}
				}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source contains a VCS configuration with an empty driver name`,
		),
		Entry(
			`duplicate VCS defaults configuration`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "test_vcs_driver" {}
					vcs "test_vcs_driver" {}
				}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source contains multiple configurations for the 'test_vcs_driver' version control system`,
		),
		Entry(
			`unrecognized VCS name`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "<unrecognized>" {}
				}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source contains configuration for an unrecognized version control system ('<unrecognized>'), the supported VCS drivers are 'test_vcs_driver'`,
		),
		Entry(
			`VCS defaults with a well-structured, but invalid body`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "test_vcs_driver" {
						unrecognized = true
					}
				}`,
			},
			`<dir>/config-0.hcl:3,7-19: Unsupported argument; An argument named "unrecognized" is not expected here.`,
		),
		Entry(
			`error normalizing driver configuration`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "test_vcs_driver" {
						filesystem_path = "~someuser/path/to/nowhere"
					}
				}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source's configuration for the 'test_vcs_driver' version control system cannot be loaded: cannot expand user-specific home dir`,
		),
	)
})
