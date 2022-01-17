package github_test

import (
	"github.com/gritcli/grit/internal/daemon/configtest"
	"github.com/gritcli/grit/internal/daemon/internal/config"
	. "github.com/gritcli/grit/internal/daemon/internal/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Config", func() {
	Describe("func DescribeSourceConfig()", func() {
		DescribeTable(
			"it describes the source",
			func(cfg Config, expect string) {
				Expect(cfg.DescribeSourceConfig()).To(Equal(expect))
			},
			Entry(
				"github.com",
				Config{Domain: "github.com"},
				"github.com",
			),
			Entry(
				"github enterprise server",
				Config{Domain: "code.example.com"},
				"code.example.com (github enterprise server)",
			),
		)
	})
})

var _ = Describe("configuration integration", func() {
	configtest.TestSourceDriver(
		SourceDriverRegistration(),
		configtest.Success(
			"authentication token",
			`source "github" "github" {
				token = "<token>"
			}`,
			Config{
				Domain: "github.com",
				Token:  "<token>",
			},
		),
		configtest.Success(
			"github enterprise server",
			`source "github" "github" {
				domain = "github.example.com"
			}`,
			Config{
				Domain: "github.example.com",
			},
		),
		configtest.Success(
			"explicit SSH key",
			`source "github" "github" {
				git {
					ssh_key {
						file = "/path/to/key"
					}
				}
			}`,
			Config{
				Domain: "github.com",
				Git: config.Git{
					SSHKeyFile: "/path/to/key",
				},
			},
		),
		configtest.Success(
			"explicit SSH key with passphrase",
			`source "github" "github" {
				git {
					ssh_key {
						file = "/path/to/key"
						passphrase = "<passphrase>"
					}
				}
			}`,
			Config{
				Domain: "github.com",
				Git: config.Git{
					SSHKeyFile:       "/path/to/key",
					SSHKeyPassphrase: "<passphrase>",
				},
			},
		),
		configtest.Success(
			"does not inherit default passphrase when SSH key is specified explicitly",
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
			Config{
				Domain: "github.com",
				Git: config.Git{
					SSHKeyFile:       "/path/to/different/key",
					SSHKeyPassphrase: "", // note: different to default git config
				},
			},
		),
		configtest.Success(
			"explicitly prefer SSH",
			`source "github" "github" {
				git {
					prefer_http = false
				}
			}`,
			Config{
				Domain: "github.com",
				Git: config.Git{
					PreferHTTP: false,
				},
			},
		),
		configtest.Success(
			"explicitly prefer HTTP",
			`source "github" "github" {
				git {
					prefer_http = true
				}
			}`,
			Config{
				Domain: "github.com",
				Git: config.Git{
					PreferHTTP: true,
				},
			},
		),
		configtest.Failure(
			`explicit SSH passphrase without key file`,
			`source "github" "github" {
				git {
					ssh_key {
						passphrase = "<passphrase>"
					}
				}
			}`,
			`<dir>/config-0.hcl:3,14-14: Missing required argument; The argument "file" is required, but no definition was found.`,
		),
	)
})
