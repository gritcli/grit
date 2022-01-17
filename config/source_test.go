package config_test

import (
	. "github.com/gritcli/grit/config"
	. "github.com/gritcli/grit/config/internal/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("func Load() (source blocks)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"explicitly enabled source",
			[]string{
				`source "test_source" "test_source_driver" {
					enabled = true
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
					VCSConfig: VCSConfigStub{Value: "<default>"},
				},
			}),
		),
		Entry(
			"explicitly disabled source",
			[]string{
				`source "test_source" "test_source_driver" {
					enabled = false
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: false,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: SourceConfigStub{
					Value:     "<default>",
					VCSConfig: VCSConfigStub{Value: "<default>"},
				},
			}),
		),
		Entry(
			"driver-specific configuration",
			[]string{
				`source "test_source" "test_source_driver" {
					value = "<explicit>"
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: SourceConfigStub{
					Value:     "<explicit>",
					VCSConfig: VCSConfigStub{Value: "<default>"},
				},
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		testLoadFailure,
		Entry(
			`empty source name`,
			[]string{
				`source "" "test_source_driver" {}`,
			},
			`<dir>/config-0.hcl: this file contains a 'source' block with an empty name`,
		),
		Entry(
			`invalid source name`,
			[]string{
				`source "<invalid>" "test_source_driver" {}`,
			},
			`<dir>/config-0.hcl: the '<invalid>' source has an invalid name, source names must contain only alpha-numeric characters and underscores`,
		),
		Entry(
			`duplicate source names`,
			[]string{
				`source "test_source" "test_source_driver" {}`,
				`source "test_source" "test_source_driver" {}`,
			},
			`<dir>/config-1.hcl: a source named 'test_source' is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`duplicate source names (case-insensitive)`,
			[]string{
				`source "test_source" "test_source_driver" {}`,
				`source "TEST_SOURCE" "test_source_driver" {}`,
			},
			`<dir>/config-1.hcl: a source named 'test_source' is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`unrecognized source driver`,
			[]string{
				`source "test_source" "<unrecognized>" {}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source uses '<unrecognized>' which is not supported, the supported drivers are: 'test_source_driver'`,
		),
		Entry(
			`source with a well-structured, but invalid body`,
			[]string{
				`source "test_source" "test_source_driver" {
					unrecognized = true
				}`,
			},
			`<dir>/config-0.hcl:2,6-18: Unsupported argument; An argument named "unrecognized" is not expected here.`,
		),
	)
})
