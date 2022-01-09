package config_test

import (
	"path/filepath"

	. "github.com/gritcli/grit/common/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
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
				withClonesDefaults(defaultConfig, Clones{
					Dir: "/path/to/clones",
				}),
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

	It("resolves the clone directory relative to the config directory", func() {
		dir, cleanup := makeConfigDir(
			`clones {
				dir = "relative/path/to/clones"
			}`,
		)
		defer cleanup()

		cfg, err := Load(dir)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg.ClonesDefaults.Dir).To(Equal(
			filepath.Join(dir, "relative/path/to/clones"),
		))
	})
})
