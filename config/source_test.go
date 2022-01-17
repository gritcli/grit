package config_test

import (
	. "github.com/gritcli/grit/config"
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
				`source "github" "github" {
					enabled = true
				}`,
			},
			defaultConfig,
		),
		Entry(
			"explicitly disabled source",
			[]string{
				`source "github" "github" {
					enabled = false
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "github",
				Enabled: false,
				Clones: Clones{
					Dir: "~/grit/github",
				},
				Driver: GitHub{
					Domain: "github.com",
				},
			}),
		),
		Entry(
			"explicit clone directory",
			[]string{
				`source "github" "github" {
					clones {
						dir = "/path/to/clones"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "github",
				Enabled: true,
				Clones: Clones{
					Dir: "/path/to/clones",
				},
				Driver: GitHub{
					Domain: "github.com",
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
				`source "" "github" {}`,
			},
			`<dir>/config-0.hcl: this file contains a 'source' block with an empty name`,
		),
		Entry(
			`invalid source name`,
			[]string{
				`source "<invalid>" "github" {}`,
			},
			`<dir>/config-0.hcl: the '<invalid>' source has an invalid name, source names must contain only alpha-numeric characters and underscores`,
		),
		Entry(
			`duplicate source names`,
			[]string{
				`source "my_company" "github" {}`,
				`source "my_company" "github" {}`,
			},
			`<dir>/config-1.hcl: a source named 'my_company' is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`duplicate source names (case-insensitive)`,
			[]string{
				`source "my_company" "github" {}`,
				`source "MY_COMPANY" "github" {}`,
			},
			`<dir>/config-1.hcl: a source named 'my_company' is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`unrecognized source driver`,
			[]string{
				`source "my_source" "<unrecognized>" {}`,
			},
			`<dir>/config-0.hcl: the 'my_source' source uses '<unrecognized>' which is not supported, the supported drivers are: 'github'`,
		),
		Entry(
			`source with a well-structured, but invalid body`,
			[]string{
				`source "github" "github" {
					unrecognized = true
				}`,
			},
			`<dir>/config-0.hcl:2,6-18: Unsupported argument; An argument named "unrecognized" is not expected here.`,
		),
	)
})
