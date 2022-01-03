package config_test

import (
	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (global git block)", func() {
	DescribeTable(
		"it returns the expected configuration",
		func(configs []string, expect Config) {
			dir, cleanup := makeConfigDir(configs...)
			defer cleanup()

			cfg, err := Load(dir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfg).To(Equal(expect))
		},
		Entry(
			"explicit private key",
			[]string{
				`git {
					private_key = "/path/to/key"
				}`,
			},
			withGlobalGit(DefaultConfig, Git{
				PrivateKey: "/path/to/key",
			}),
		),
		Entry(
			"explicitly prefer SSH",
			[]string{
				`git {
					prefer_http = false
				}`,
			},
			DefaultConfig,
		),
		Entry(
			"explicitly prefer HTTP",
			[]string{
				`git {
					prefer_http = true
				}`,
			},
			withGlobalGit(DefaultConfig, Git{
				PreferHTTP: true,
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		func(configs []string, expect string) {
			dir, cleanup := makeConfigDir(configs...)
			defer cleanup()

			_, err := Load(dir)
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(expect)), err.Error())
		},
		Entry(
			`multiple files with git blocks`,
			[]string{`git {}`, `git {}`},
			`the global git configuration has already been defined in`,
		),
	)
})

// withGlobalGit returns a copy of cfg with a different git configuration.
func withGlobalGit(cfg Config, g Git) Config {
	cfg.GlobalGit = g
	return cfg
}
