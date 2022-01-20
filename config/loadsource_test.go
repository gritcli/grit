package config_test

import (
	. "github.com/gritcli/grit/config"
	. "github.com/gritcli/grit/config/internal/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("func Load() (source configuration)", func() {
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
			`<dir>/config-0.hcl: source configurations must provide a name`,
		),
		Entry(
			`invalid source name`,
			[]string{
				`source "<invalid>" "test_source_driver" {}`,
			},
			`<dir>/config-0.hcl: '<invalid>' is not a valid source name, valid characters are ASCII letters, numbers and underscore`,
		),
		Entry(
			`duplicate source names`,
			[]string{
				`source "test_source" "test_source_driver" {}`,
				`source "test_source" "test_source_driver" {}`,
			},
			`<dir>/config-1.hcl: the 'test_source' source conflicts with a source of the same name in <dir>/config-0.hcl (source names are case-insensitive)`,
		),
		Entry(
			`duplicate source names (case-insensitive)`,
			[]string{
				`source "test_source" "test_source_driver" {}`,
				`source "TEST_SOURCE" "test_source_driver" {}`,
			},
			`<dir>/config-1.hcl: the 'TEST_SOURCE' source conflicts with a source of the same name in <dir>/config-0.hcl (source names are case-insensitive)`,
		),
		Entry(
			`empty driver name`,
			[]string{
				`source "test_source" "" {}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source has an empty driver name`,
		),
		Entry(
			`unrecognized source driver name`,
			[]string{
				`source "test_source" "<unrecognized>" {}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source uses an unrecognized driver ('<unrecognized>'), the supported source drivers are 'test_source_driver'`,
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
