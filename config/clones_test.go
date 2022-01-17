package config_test

import (
	. "github.com/gritcli/grit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("func Load() (global clones block)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"explicit clone directory",
			[]string{
				`clones {
					dir = "/path/to/clones"
				}`,
			},
			withSource(
				defaultConfig,
				Source{
					Name:    "github",
					Enabled: true,
					Clones: Clones{ // base directory inherited from the clones default block
						Dir: "/path/to/clones/github",
					},
					Driver: GitHub{
						Domain: "github.com",
					},
				},
			),
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
			`<dir>/config-1.hcl: a 'clones' defaults block is already defined in <dir>/config-0.hcl`,
		),
	)
})
