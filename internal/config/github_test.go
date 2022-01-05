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
		Entry(
			"explicit private key",
			[]string{
				`source "github" "github" {
					git {
						private_key = "/path/to/key"
					}
				}`,
			},
			withSource(
				defaultConfig,
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHubConfig{
						Domain: "github.com",
						Git: Git{
							PrivateKey: "/path/to/key",
						},
					},
				},
			),
		),
		Entry(
			"explicitly prefer SSH",
			[]string{
				`source "github" "github" {
					git {
						prefer_http = false
					}
				}`,
			},
			defaultConfig,
		),
		Entry(
			"explicitly prefer HTTP",
			[]string{
				`source "github" "github" {
					git {
						prefer_http = true
					}
				}`,
			},
			withSource(
				defaultConfig,
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHubConfig{
						Domain: "github.com",
						Git: Git{
							PreferHTTP: true,
						},
					},
				},
			),
		),
	)
})
