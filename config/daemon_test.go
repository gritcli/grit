package config_test

import (
	"path/filepath"

	. "github.com/gritcli/grit/config"
	"github.com/gritcli/grit/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (daemon block)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"non-standard daemon socket",
			[]string{
				`daemon {
					socket = "/path/to/socket"
				}`,
			},
			withDaemon(defaultConfig, Daemon{
				Socket: "/path/to/socket",
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		testLoadFailure,
		Entry(
			`multiple files with daemon blocks`,
			[]string{
				`daemon {}`,
				`daemon {}`,
			},
			`<dir>/config-1.hcl: a 'daemon' block is already defined in <dir>/config-0.hcl`,
		),
	)

	It("resolves the socket path relative to the config directory", func() {
		dir, cleanup := makeConfigDir(
			`daemon {
				socket = "relative/path/to/socket"
			}`,
		)
		defer cleanup()

		// TODO: don't test using built-ins
		cfg, err := Load(dir, &registry.Registry{
			Parent: &registry.BuiltIns,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg.Daemon.Socket).To(Equal(
			filepath.Join(dir, "relative/path/to/socket"),
		))
	})
})
