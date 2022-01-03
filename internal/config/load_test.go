package config_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (github source)", func() {
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
		Entry(
			"empty directory is equivalent to the default",
			[]string{},
			DefaultConfig,
		),
		Entry(
			"empty file is equivalent to the default",
			[]string{``},
			DefaultConfig,
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
			`syntax error`,
			[]string{`<invalid>`},
			`Argument or block definition required; An argument or block definition is required here.`,
		),
		Entry(
			`multiple files with daemon blocks`,
			[]string{`daemon {}`, `daemon {}`},
			`the daemon configuration has already been defined in`,
		),
		Entry(
			`duplicate source names`,
			[]string{`source "my_company" "github" {}`, `source "my_company" "github" {}`},
			`the 'my_company' repository source has already been defined in`,
		),
		Entry(
			`empty source name`,
			[]string{`source "" "github" {}`},
			`the '' repository source is invalid: source name must not be empty`,
		),
		Entry(
			`invalid source name`,
			[]string{`source "<invalid>" "github" {}`},
			`the '<invalid>' repository source is invalid: source name must contain only alpha-numeric characters and underscores`,
		),
	)

	It("returns the default configuration when passed a non-existant directory", func() {
		cfg, err := Load("./does-not-exist")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg).To(Equal(DefaultConfig))
	})

	It("ignores non-HCL files, directories and HCL files that begin with an underscore", func() {
		dir, err := os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())
		defer os.RemoveAll(dir)

		err = os.Mkdir(filepath.Join(dir, "subdirectory"), 0700)
		Expect(err).ShouldNot(HaveOccurred())

		err = os.WriteFile(filepath.Join(dir, "subdirectory", "should-be-ignored.txt"), []byte("<invalid config>"), 0600)
		Expect(err).ShouldNot(HaveOccurred())

		err = os.WriteFile(filepath.Join(dir, "not-hcl.txt"), []byte("<invalid config>"), 0600)
		Expect(err).ShouldNot(HaveOccurred())

		err = os.WriteFile(filepath.Join(dir, "_underscore.hcl"), []byte("<invalid config>"), 0600)
		Expect(err).ShouldNot(HaveOccurred())

		cfg, err := Load(dir)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg).To(Equal(DefaultConfig))
	})
})

// makeConfigDir makes a temporary config directory containing config files
// containing the given configuration content.
func makeConfigDir(configs ...string) (dir string, cleanup func()) {
	dir, err := os.MkdirTemp("", "")
	Expect(err).ShouldNot(HaveOccurred())

	for i, cfg := range configs {
		err := os.WriteFile(
			filepath.Join(dir, fmt.Sprintf("config-%d.hcl", i)),
			[]byte(cfg),
			0600,
		)
		Expect(err).ShouldNot(HaveOccurred())
	}

	return dir, func() {
		os.RemoveAll(dir)
	}
}

// withDaemon returns a copy of cfg with a different daemon configuration.
func withDaemon(cfg Config, d Daemon) Config {
	cfg.Daemon = d
	return cfg
}

// withSource returns a copy of cfg with an additional repository source.
func withSource(cfg Config, src Source) Config {
	prev := cfg.Sources
	cfg.Sources = map[string]Source{}

	for n, s := range prev {
		cfg.Sources[n] = s
	}

	cfg.Sources[src.Name] = src

	return cfg
}
