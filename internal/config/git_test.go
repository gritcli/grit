package config_test

import (
	"path/filepath"

	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (global git block)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"explicit private key",
			[]string{
				`git {
					private_key = "/path/to/key"
				}`,
			},
			withSource(
				withGlobalGit(defaultConfig, Git{
					PrivateKey: "/path/to/key",
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHub{
						Domain: "github.com",
						Git: Git{ // inherited from global git block
							PrivateKey: "/path/to/key",
						},
					},
				},
			),
		),
		Entry(
			"explicit private key with passphrase",
			[]string{
				`git {
					private_key = "/path/to/key"
					passphrase = "<passphrase>"
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
						Git: Git{ // inherited from global git block
							PrivateKey: "/path/to/key",
							Passphrase: "<passphrase>",
						},
					},
				},
			),
		),
		Entry(
			"explicitly prefer SSH",
			[]string{
				`git {
					prefer_http = false
				}`,
			},
			defaultConfig,
		),
		Entry(
			"explicitly prefer HTTP",
			[]string{
				`git {
					prefer_http = true
				}`,
			},
			withSource(
				withGlobalGit(defaultConfig, Git{
					PreferHTTP: true,
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHub{
						Domain: "github.com",
						Git: Git{ // inherited from global git block
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
			`multiple files with global git blocks`,
			[]string{
				`git {}`,
				`git {}`,
			},
			`<dir>/config-1.hcl: a global 'git' block is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`explicit passphrase without private key`,
			[]string{
				`git {
					passphrase = "<passphrase>"
				}`,
			},
			`<dir>/config-0.hcl: the global 'git' block is invalid: passphrase present without specifying a private key file`,
		),
	)

	It("resolves the private key path relative to the config directory", func() {
		dir, cleanup := makeConfigDir(
			`git {
				private_key = "relative/path/to/key"
			}`,
		)
		defer cleanup()

		cfg, err := Load(dir)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg.GlobalGit.PrivateKey).To(Equal(
			filepath.Join(dir, "relative/path/to/key"),
		))
	})
})
