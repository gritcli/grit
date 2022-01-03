package config_test

import (
	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (source blocks)", func() {
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
			"explicitly enabled source",
			[]string{
				`source "github" "github" {
					enabled = true
				}`,
			},
			DefaultConfig,
		),
		Entry(
			"explicitly disabled source",
			[]string{
				`source "github" "github" {
					enabled = false
				}`,
			},
			withSource(DefaultConfig, Source{
				Name:    "github",
				Enabled: false,
				Config: GitHubConfig{
					Domain: "github.com",
				},
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
			`empty source name`,
			[]string{`source "" "github" {}`},
			`repository sources must not have empty names`,
		),
		Entry(
			`invalid source name`,
			[]string{`source "<invalid>" "github" {}`},
			`the '<invalid>' repository source has an invalid name, names must contain only alpha-numeric characters and underscores`,
		),
		Entry(
			`unrecognized source implementation`,
			[]string{`source "my_source" "<unrecognized>" {}`},
			`'<unrecognized>' is not a recognized repository source implementation, expected 'github'`,
		),
		Entry(
			`duplicate source names`,
			[]string{`source "my_company" "github" {}`, `source "my_company" "github" {}`},
			`a repository source named 'my_company' has already been defined in`,
		),
	)
})

// withSource returns a copy of cfg with an additional repository source.
func withSource(cfg Config, src Source) Config {
	prev := cfg.Sources
	cfg.Sources = map[string]Source{}

	for n, s := range prev {
		cfg.Sources[n] = s
	}

	cfg.Sources[src.Name] = src

	return cfg
}
