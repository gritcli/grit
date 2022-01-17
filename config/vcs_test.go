package config_test

import (
	. "github.com/gritcli/grit/config"
	. "github.com/gritcli/grit/config/internal/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("func Load() (vcs blocks)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"sources inherit the default VCS configuration",
			[]string{
				`vcs "test_vcs_driver" {
					value = "<explicit>"
				}

				source "test_source" "test_source_driver" {}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: SourceConfigStub{
					Value:     "<default>",
					VCSConfig: VCSConfigStub{Value: "<explicit>"},
				},
			}),
		),
		Entry(
			"sources can specify a VCS configuration",
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "test_vcs_driver" {
						value = "<explicit>"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: SourceConfigStub{
					Value:     "<default>",
					VCSConfig: VCSConfigStub{Value: "<default><explicit>"},
				},
			}),
		),
		Entry(
			"sources can override the default VCS configuration",
			[]string{
				`vcs "test_vcs_driver" {
					value = "<explicit default>"
				}

				source "test_source" "test_source_driver" {
					vcs "test_vcs_driver" {
						value = "<override>"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: SourceConfigStub{
					Value:     "<default>",
					VCSConfig: VCSConfigStub{Value: "<explicit default><override>"},
				},
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the global configuration",
		testLoadFailure,
		Entry(
			`empty VCS driver alias`,
			[]string{
				`vcs "" {}`,
			},
			`<dir>/config-0.hcl: this file contains a 'vcs' block with an empty driver alias`,
		),
		Entry(
			`duplicate VCS defaults configuration`,
			[]string{
				`vcs "test_vcs_driver" {}`,
				`vcs "test_vcs_driver" {}`,
			},
			`<dir>/config-1.hcl: defaults for the 'test_vcs_driver' version control system are already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`unrecognized VCS alias`,
			[]string{
				`vcs "<unrecognized>" {}`,
			},
			`<dir>/config-0.hcl: the '<unrecognized>' version control system is not supported, the supported drivers are: 'test_vcs_driver'`,
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
	)

	DescribeTable(
		"it returns an error if there is a problem with a source-specific configuration",
		testLoadFailure,
		Entry(
			`empty VCS driver alias`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "" {}
				}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source contains a 'vcs' block with an empty driver alias`,
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
			`unrecognized VCS alias`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "<unrecognized>" {}
				}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source contains configuration for the '<unrecognized>' version control system, which is not supported, the supported drivers are: 'test_vcs_driver'`,
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
	)
})
