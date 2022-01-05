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
			"explicit SSH key",
			[]string{
				`git {
					ssh_key {
						file = "/path/to/key"
					}
				}`,
			},
			withSource(
				withGlobalGit(defaultConfig, Git{
					SSHKeyFile: "/path/to/key",
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHub{
						Domain: "github.com",
						Git: Git{ // inherited from global git block
							SSHKeyFile: "/path/to/key",
						},
					},
				},
			),
		),
		Entry(
			"explicit SSH key with passphrase",
			[]string{
				`git {
					ssh_key {
						file = "/path/to/key"
						passphrase = "<passphrase>"
					}
				}`,
			},
			withSource(
				withGlobalGit(defaultConfig, Git{
					SSHKeyFile:       "/path/to/key",
					SSHKeyPassphrase: "<passphrase>",
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Config: GitHub{
						Domain: "github.com",
						Git: Git{ // inherited from global git block
							SSHKeyFile:       "/path/to/key",
							SSHKeyPassphrase: "<passphrase>",
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
			`explicit SSH passphrase without key file`,
			[]string{
				`git {
					ssh_key {
						passphrase = "<passphrase>"
					}
				}`,
			},
			`<dir>/config-0.hcl:2,14-14: Missing required argument; The argument "file" is required, but no definition was found.`,
		),
	)

	It("resolves the SSH key path relative to the config directory", func() {
		dir, cleanup := makeConfigDir(
			`git {
				ssh_key {
					file = "relative/path/to/key"
				}
			}`,
		)
		defer cleanup()

		cfg, err := Load(dir)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg.GlobalGit.SSHKeyFile).To(Equal(
			filepath.Join(dir, "relative/path/to/key"),
		))
	})
})
