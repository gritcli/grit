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
				Clones: Clones{
					Dir: "~/grit",
				},
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
				Clones: Clones{
					Dir: "~/grit",
				},
				Config: GitHub{
					Domain: "github.example.com",
				},
			}),
		),
		Entry(
			"explicit SSH key",
			[]string{
				`source "github" "github" {
					git {
						ssh_key {
							file = "/path/to/key"
						}
					}
				}`,
			},
			withSource(
				defaultConfig,
				Source{
					Name:    "github",
					Enabled: true,
					Clones: Clones{
						Dir: "~/grit",
					},
					Config: GitHub{
						Domain: "github.com",
						Git: Git{
							SSHKeyFile: "/path/to/key",
						},
					},
				},
			),
		),
		Entry(
			"explicit SSH key with passphrase",
			[]string{
				`source "github" "github" {
					git {
						ssh_key {
							file = "/path/to/key"
							passphrase = "<passphrase>"
						}
					}
				}`,
			},
			withSource(
				defaultConfig,
				Source{
					Name:    "github",
					Enabled: true,
					Clones: Clones{
						Dir: "~/grit",
					},
					Config: GitHub{
						Domain: "github.com",
						Git: Git{
							SSHKeyFile:       "/path/to/key",
							SSHKeyPassphrase: "<passphrase>",
						},
					},
				},
			),
		),
		Entry(
			"does not inherit default passphrase when SSH key is specified explicitly",
			[]string{
				`git {
					ssh_key {
						file = "/path/to/key"
						passphrase = "<passphrase>"
					}
				}

				source "github" "github" {
					git {
						ssh_key {
							file = "/path/to/different/key"
						}
					}
				}`,
			},
			withSource(
				withGitDefaults(defaultConfig, Git{
					SSHKeyFile:       "/path/to/key",
					SSHKeyPassphrase: "<passphrase>",
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Clones: Clones{
						Dir: "~/grit",
					},
					Config: GitHub{
						Domain: "github.com",
						Git: Git{
							SSHKeyFile:       "/path/to/different/key",
							SSHKeyPassphrase: "", // note: different to default git config
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
					Clones: Clones{
						Dir: "~/grit",
					},
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
			`explicit SSH passphrase without key file`,
			[]string{
				`source "github" "github" {
					git {
						ssh_key {
							passphrase = "<passphrase>"
						}
					}
				}`,
			},
			`<dir>/config-0.hcl:3,15-15: Missing required argument; The argument "file" is required, but no definition was found.`,
		),
	)
})
