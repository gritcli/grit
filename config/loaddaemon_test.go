package config_test

import (
	"path/filepath"

	. "github.com/gritcli/grit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (daemon configuration)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"explicit daemon socket",
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
			`<dir>/config-1.hcl: the daemon configuration is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`unexpandable daemon socket path`,
			[]string{
				`daemon {
					socket = "~someuser/path/to/socket"
				}`,
			},
			`<dir>/config-0.hcl: unable to resolve daemon socket path: cannot expand user-specific home dir (~someuser/path/to/socket)`,
		),
	)

	Context("when the default daemon socket can not be resolved", func() {
		var original string

		BeforeEach(func() {
			// HACK: We really shouldn't manipulate (or even have) global
			// variables like this, but it's the only cross-platform way to
			// force the home directory resolution to fail.
			original = DefaultDaemonSocket
			DefaultDaemonSocket = "~someuser/path/to/socket"
		})

		AfterEach(func() {
			DefaultDaemonSocket = original
		})

		DescribeTable(
			"it returns an error",
			testLoadFailure,
			Entry(
				`unexpandable default daemon socket`,
				[]string{},
				`unable to resolve default daemon socket path: cannot expand user-specific home dir (~someuser/path/to/socket)`,
			),
		)
	})

	It("resolves the socket path relative to the config directory", func() {
		dir, cleanup := makeConfigDir(
			`daemon {
				socket = "relative/path/to/socket"
			}`,
		)
		defer cleanup()

		cfg, err := Load(dir, nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg.Daemon.Socket).To(Equal(
			filepath.Join(dir, "relative/path/to/socket"),
		))
	})
})
