package config_test

import (
	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (github source)", func() {
	DescribeTable(
		"it returns the expected configuration",
		func(config string, expect Config) {
			dir, cleanup := makeConfigDir(config)
			defer cleanup()

			cfg, err := Load(dir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfg).To(Equal(expect))
		},
		Entry(
			"token",
			`source "github" "github" {
				token = "<token>"
			}`,
			withSource(DefaultConfig, Source{
				Name:    "github",
				Enabled: true,
				Config: GitHubConfig{
					Domain: "github.com",
					Token:  "<token>",
				},
			}),
		),
		Entry(
			"github enterprise",
			`source "my_company" "github" {
				domain = "github.example.com"
			}`,
			withSource(DefaultConfig, Source{
				Name:    "my_company",
				Enabled: true,
				Config: GitHubConfig{
					Domain: "github.example.com",
				},
			}),
		),
	)
})
