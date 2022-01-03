package config_test

import (
	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (daemon block)", func() {
	DescribeTable(
		"it returns the expected configuration",
		func(configs []string, expect Config) {
			dir, cleanup := makeConfigDir(configs...)
			defer cleanup()

			cfg, err := Load(dir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfg).To(Equal(expect))
		},
		Entry(
			"non-standard daemon socket",
			[]string{
				`daemon {
					socket = "/path/to/socket"
				}`,
			},
			withDaemon(DefaultConfig, Daemon{
				Socket: "/path/to/socket",
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		func(configs []string, expect string) {
			dir, cleanup := makeConfigDir(configs...)
			defer cleanup()

			_, err := Load(dir)
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(expect)), err.Error())
		},
		Entry(
			`multiple files with daemon blocks`,
			[]string{`daemon {}`, `daemon {}`},
			`the daemon configuration has already been defined in`,
		),
	)
})
