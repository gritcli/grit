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
				Config: GitHub{
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
				Config: GitHub{
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
					Config: GitHub{
						Domain: "github.com",
						Git: Git{
							PrivateKey: "/path/to/key",
						},
					},
				},
			),
		),
		Entry(
			"explicit private key with passphrase",
			[]string{
				`source "github" "github" {
					git {
						private_key = "/path/to/key"
						passphrase = "<passphrase>"
					}
				}`,
			},
			withSource(
				defaultConfig,
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHub{
						Domain: "github.com",
						Git: Git{
							PrivateKey: "/path/to/key",
							Passphrase: "<passphrase>",
						},
					},
				},
			),
		),
		Entry(
			"does not inherit global passphase when private key is specified explicitly",
			[]string{
				`git {
					private_key = "/path/to/key"
					passphrase = "<passphrase>"
				}

				source "github" "github" {
					git {
						private_key = "/path/to/different/key"
					}
				}`,
			},
			withSource(
				withGlobalGit(defaultConfig, Git{
					PrivateKey: "/path/to/key",
					Passphrase: "<passphrase>",
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHub{
						Domain: "github.com",
						Git: Git{
							PrivateKey: "/path/to/different/key",
							Passphrase: "", // note: different to global git config
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
					Config: GitHub{
						Domain: "github.com",
						Git: Git{
							PreferHTTP: true,
						},
					},
				},
			),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		testLoadFailure,
		Entry(
			`explicit passphrase without private key`,
			[]string{
				`source "github" "github" {
					git {
						passphrase = "<passphrase>"
					}
				}`,
			},
			`<dir>/config-0.hcl: the 'github' repository source is invalid: the 'git' block is invalid: passphrase present without specifying a private key file`,
		),
	)
})
