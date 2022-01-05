package config_test

import (
	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("func Load() (github source)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"authentication token",
			[]string{
				`source "github" "github" {
					token = "<token>"
				}`,
			},
			withSource(defaultConfig, Source{
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
			[]string{
				`source "my_company" "github" {
					domain = "github.example.com"
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "my_company",
				Enabled: true,
				Config: GitHubConfig{
					Domain: "github.example.com",
				},
			}),
		),
	)
})
