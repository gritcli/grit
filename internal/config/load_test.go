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

var _ = Describe("func Load()", func() {
	DescribeTable(
		"it returns the expected configuration",
		func(dir string, expect Config) {
			cfg, err := Load(dir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfg).To(Equal(expect))
		},
		Entry(
			"default configuration",
			"testdata/valid/default",
			DefaultConfig,
		),
		Entry(
			"empty configuration file (should be the same as the default)",
			"testdata/valid/empty-file",
			DefaultConfig,
		),
		Entry(
			"empty configuration directory (should be the same as the default)",
			"testdata/valid/empty-dir",
			DefaultConfig,
		),
		Entry(
			`not existent directory (should be the same as the default)`,
			`testdata/valid/non-existent`,
			DefaultConfig,
		),
		Entry(
			"ignores non-HCL files, directories and files beginning with underscores",
			"testdata/valid/ignore",
			DefaultConfig,
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		func(dir string, expect string) {
			_, err := Load(dir)
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expect), err.Error())
		},
		Entry(
			`syntax error`,
			`testdata/invalid/syntax-error`,
			`testdata/invalid/syntax-error/grit.hcl:1,1-2: Argument or block definition required; An argument or block definition is required here.`,
		),
		Entry(
			`multiple files with daemon blocks`,
			`testdata/invalid/multiple-files-with-daemon-block`,
			`testdata/invalid/multiple-files-with-daemon-block/b.hcl: the daemon configuration has already been defined in testdata/invalid/multiple-files-with-daemon-block/a.hcl`,
		),
		Entry(
			`duplicate source names`,
			`testdata/invalid/duplicate-source-names`,
			`testdata/invalid/duplicate-source-names/b.hcl: the 'my-company' repository source has already been defined in testdata/invalid/duplicate-source-names/a.hcl`,
		),
	)
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
