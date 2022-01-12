package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/gritcli/grit/internal/daemon/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var _ = Describe("func Load()", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"empty directory is equivalent to the default",
			[]string{},
			defaultConfig,
		),
		Entry(
			"empty file is equivalent to the default",
			[]string{
				``,
			},
			defaultConfig,
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		testLoadFailure,
		Entry(
			`syntax error`,
			[]string{
				`<invalid>`,
			},
			`<dir>/config-0.hcl:1,1-2: Argument or block definition required; An argument or block definition is required here.`,
		),
	)

	It("returns the default configuration when passed a non-existent directory", func() {
		cfg, err := Load("./does-not-exist")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg).To(Equal(defaultConfig))
	})

	It("ignores non-HCL files, directories and HCL files that begin with a dot or an underscore", func() {
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

		err = os.WriteFile(filepath.Join(dir, ".dot.hcl"), []byte("<invalid config>"), 0600)
		Expect(err).ShouldNot(HaveOccurred())

		cfg, err := Load(dir)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg).To(Equal(defaultConfig))
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

// testLoadSuccess is a function for use in DescribeTable that tests for success
// cases when loading configuration files.
func testLoadSuccess(configs []string, expect Config) {
	dir, cleanup := makeConfigDir(configs...)
	defer cleanup()

	cfg, err := Load(dir)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(cfg).To(Equal(expect))
}

// testLoadFailure is a function for use in DescribeTable that tests for failure
// cases when loading configuration files.
//
// The text "<dir>" can be used in the expected message as a placeholder for the
// actual temporary directory used during the test.
func testLoadFailure(configs []string, expect string) {
	orig := format.TruncatedDiff
	format.TruncatedDiff = false
	defer func() {
		format.TruncatedDiff = orig
	}()

	dir, cleanup := makeConfigDir(configs...)
	defer cleanup()

	_, err := Load(dir)
	Expect(err).Should(HaveOccurred())

	message := strings.ReplaceAll(err.Error(), dir, "<dir>")
	Expect(message).To(Equal(expect))
}
