package config_test

import (
	. "github.com/gritcli/grit/config"
	. "github.com/gritcli/grit/config/internal/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("func Load() (clones configuration)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"sources use a directory within the default clone directory by default",
			[]string{
				`clones {
					dir = "/path/to/clones"
				}

				source "test_source" "test_source_driver" {}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "/path/to/clones/test_source",
				},
				Driver: SourceConfigStub{
					Value:     "<default>",
					VCSConfig: VCSConfigStub{Value: "<default>"},
				},
			}),
		),
		Entry(
			"sources can specify a clones configuration",
			[]string{
				`source "test_source" "test_source_driver" {
					clones {
						dir = "/path/to/clones"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "/path/to/clones",
				},
				Driver: SourceConfigStub{
					Value:     "<default>",
					VCSConfig: VCSConfigStub{Value: "<default>"},
				},
			}),
		),
		Entry(
			"sources can override the default clone directory",
			[]string{
				`clones {
					dir = "/path/to/clones"
				}

				source "test_source" "test_source_driver" {
					clones {
						dir = "/path/to/elsewhere"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "/path/to/elsewhere",
				},
				Driver: SourceConfigStub{
					Value:     "<default>",
					VCSConfig: VCSConfigStub{Value: "<default>"},
				},
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		testLoadFailure,
		Entry(
			`multiple files with clones defaults blocks`,
			[]string{
				`clones {}`,
				`clones {}`,
			},
			`<dir>/config-1.hcl: the global clones configuration is already defined in <dir>/config-0.hcl`,
		),
	)
})
