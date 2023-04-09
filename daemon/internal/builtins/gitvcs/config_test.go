package gitvcs_test

import (
	. "github.com/gritcli/grit/daemon/internal/builtins/gitvcs"
	"github.com/gritcli/grit/driver/configtest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Config", func() {
	Describe("func DescribeVCSConfig()", func() {
		DescribeTable(
			"it describes the source",
			func(cfg Config, expect string) {
				Expect(cfg.DescribeVCSConfig()).To(Equal(expect))
			},
			Entry(
				"default",
				Config{},
				"use ssh agent",
			),
			Entry(
				"explicit key",
				Config{
					SSHKeyFile: "/path/to/key.pem",
				},
				"use ssh key (key.pem)",
			),
			Entry(
				"prefer HTTP",
				Config{
					PreferHTTP: true,
				},
				"use ssh agent, prefer http",
			),
		)
	})
})

var _ = Describe("type configLoader", func() {
	configtest.TestVCSDriver(
		Registration,
		Config{},
		configtest.VCSSuccess(
			"explicit SSH key",
			`vcs "git" {
				ssh_key {
					file = "/path/to/key"
				}
			}`,
			Config{
				SSHKeyFile: "/path/to/key",
			},
		),
		configtest.VCSSuccess(
			"explicit SSH key with passphrase",
			`vcs "git" {
				ssh_key {
					file = "/path/to/key"
					passphrase = "<passphrase>"
				}
			}`,
			Config{
				SSHKeyFile:       "/path/to/key",
				SSHKeyPassphrase: "<passphrase>",
			},
		),
		configtest.VCSSuccess(
			"explicitly prefer SSH",
			`vcs "git" {
				prefer_http = false
			}`,
			Config{
				PreferHTTP: false,
			},
		),
		configtest.VCSSuccess(
			"explicitly prefer HTTP",
			`vcs "git" {
				prefer_http = true
			}`,
			Config{
				PreferHTTP: true,
			},
		),
		configtest.VCSFailure(
			`explicit SSH passphrase without key file`,
			`vcs "git" {
				ssh_key {
					passphrase = "<passphrase>"
				}
			}`,
			`<dir>/config-0.hcl:2,13-13: Missing required argument; The argument "file" is required, but no definition was found.`,
		),
		configtest.VCSSourceSpecificSuccess(
			"override SSH key",
			`vcs "git" {
				ssh_key {
					file = "/path/to/key"
					passphrase = "<passphrase>"
				}
			}`,
			`vcs "git" {
				ssh_key {
					file = "/path/to/override"
				}
			}`,
			Config{
				SSHKeyFile: "/path/to/override",
				// Note, passphrase is not inherited from default.
			},
		),
		configtest.VCSSourceSpecificSuccess(
			"override SSH key with passphrase",
			`vcs "git" {
				ssh_key {
					file = "/path/to/key"
					passphrase = "<passphrase>"
				}
			}`,
			`vcs "git" {
				ssh_key {
					file = "/path/to/override"
					passphrase = "<override>"
				}
			}`,
			Config{
				SSHKeyFile:       "/path/to/override",
				SSHKeyPassphrase: "<override>",
			},
		),
		configtest.VCSSourceSpecificSuccess(
			"override prefer HTTP",
			`vcs "git" {
				prefer_http = true
			}`,
			`vcs "git" {
				prefer_http = false
			}`,
			Config{
				PreferHTTP: false,
			},
		),
	)

})
