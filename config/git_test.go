package config_test

import (
	"path/filepath"

	. "github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/sourcedriver/githubsource"
	"github.com/gritcli/grit/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// TODO: don't test using built-ins
type GitHub = githubsource.Config

var _ = Describe("func Load() (git defaults block)", func() {
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
				withGitDefaults(defaultConfig, Git{
					SSHKeyFile: "/path/to/key",
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Clones: Clones{
						Dir: "~/grit/github",
					},
					Driver: GitHub{
						Domain: "github.com",
						Git: Git{ // inherited from git defaults block
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
				withGitDefaults(defaultConfig, Git{
					SSHKeyFile:       "/path/to/key",
					SSHKeyPassphrase: "<passphrase>",
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Clones: Clones{
						Dir: "~/grit/github",
					},
					Driver: GitHub{
						Domain: "github.com",
						Git: Git{ // inherited from git defaults block
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
				withGitDefaults(defaultConfig, Git{
					PreferHTTP: true,
				}),
				Source{
					Name:    "github",
					Enabled: true,
					Clones: Clones{
						Dir: "~/grit/github",
					},
					Driver: GitHub{
						Domain: "github.com",
						Git: Git{ // inherited from git defaults block
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
			`multiple files with git defaults blocks`,
			[]string{
				`git {}`,
				`git {}`,
			},
			`<dir>/config-1.hcl: a 'git' defaults block is already defined in <dir>/config-0.hcl`,
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

		// TODO: don't test using built-ins
		cfg, err := Load(dir, &registry.Registry{
			Parent: &registry.BuiltIns,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg.GitDefaults.SSHKeyFile).To(Equal(
			filepath.Join(dir, "relative/path/to/key"),
		))
	})
})
