package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/gritcli/grit/internal/stubs"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	. "github.com/onsi/ginkgo/v2"
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
		cfg, err := Load("./does-not-exist", nil)
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

		cfg, err := Load(dir, nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cfg).To(Equal(defaultConfig))
	})

	It("uses the DefaultDirectory by default", func() {
		// HACK: We really shouldn't manipulate (or even have) global variables
		// like this.
		original := DefaultDirectory
		DefaultDirectory = "~someuser/path/to/config" // force a failure so we know this path was used
		defer func() { DefaultDirectory = original }()

		_, err := Load("", nil)
		Expect(err).To(MatchError("unable to resolve configuration directory: cannot expand user-specific home dir (~someuser/path/to/config)"))
	})

	It("returns an error if the config directory cannot be resolved", func() {
		_, err := Load("~someuser/path/to/config", nil)
		Expect(err).To(MatchError("unable to resolve configuration directory: cannot expand user-specific home dir (~someuser/path/to/config)"))
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
func testLoadSuccess(
	configs []string,
	expect Config,
	hooks ...func(r *DriverRegistry),
) {
	dir, cleanup := makeConfigDir(configs...)
	defer cleanup()

	r := newRegistry()
	for _, h := range hooks {
		h(r)
	}

	cfg, err := Load(dir, r)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(cfg).To(Equal(expect))
}

// testLoadFailure is a function for use in DescribeTable that tests for failure
// cases when loading configuration files.
//
// The text "<dir>" can be used in the expected message as a placeholder for the
// actual temporary directory used during the test.
func testLoadFailure(
	configs []string,
	expect string,
	hooks ...func(r *DriverRegistry),
) {
	orig := format.TruncatedDiff
	format.TruncatedDiff = false
	defer func() {
		format.TruncatedDiff = orig
	}()

	dir, cleanup := makeConfigDir(configs...)
	defer cleanup()

	r := newRegistry()
	for _, h := range hooks {
		h(r)
	}

	_, err := Load(dir, r)
	Expect(err).Should(HaveOccurred())

	message := strings.ReplaceAll(err.Error(), dir, "<dir>")
	Expect(message).To(Equal(expect))
}

const (
	testSourceDriverName = "test_source_driver"
	testVCSDriverName    = "test_vcs_driver"
)

// newRegistry returns the registry to use for Load() tests.
func newRegistry() *DriverRegistry {
	reg := &DriverRegistry{}

	reg.RegisterSourceDriver(
		testSourceDriverName,
		sourcedriver.Registration{
			Name:         testSourceDriverName,
			ConfigLoader: newSourceLoader(),
		},
	)

	reg.RegisterVCSDriver(
		testVCSDriverName,
		vcsdriver.Registration{
			Name:         testVCSDriverName,
			ConfigLoader: newVCSLoader(),
		},
	)

	return reg
}

func newSourceLoader() *stubs.SourceConfigLoader {
	return &stubs.SourceConfigLoader{
		UnmarshalFunc: func(
			ctx sourcedriver.ConfigContext,
			b hcl.Body,
		) (sourcedriver.Config, error) {
			var s stubs.SourceConfigSchema
			if diags := gohcl.DecodeBody(b, ctx.EvalContext(), &s); diags.HasErrors() {
				return nil, diags
			}

			cfg := &stubs.SourceConfig{
				ArbitraryAttribute: "<default>",
			}

			if s.ArbitraryAttribute != "" {
				cfg.ArbitraryAttribute = s.ArbitraryAttribute
			}

			if s.FilesystemPath != "" {
				cfg.FilesystemPath = s.FilesystemPath
			}

			if err := ctx.NormalizePath(&cfg.FilesystemPath); err != nil {
				return nil, err
			}

			vcsConfig := &stubs.VCSConfig{}
			if err := ctx.UnmarshalVCSConfig(testVCSDriverName, &vcsConfig); err != nil {
				return nil, err
			}

			cfg.VCSs = map[string]vcsdriver.Config{
				testVCSDriverName: vcsConfig,
			}

			return cfg, nil
		},
	}
}

func newVCSLoader() *stubs.VCSConfigLoader {
	return &stubs.VCSConfigLoader{
		DefaultsFunc: func(
			vcsdriver.ConfigContext,
		) (vcsdriver.Config, error) {
			return &stubs.VCSConfig{
				ArbitraryAttribute: "<default>",
			}, nil
		},
		UnmarshalAndMergeFunc: func(
			ctx vcsdriver.ConfigContext,
			c vcsdriver.Config,
			b hcl.Body,
		) (vcsdriver.Config, error) {
			var s stubs.VCSConfigSchema
			if diags := gohcl.DecodeBody(b, ctx.EvalContext(), &s); diags.HasErrors() {
				return nil, diags
			}

			cfg := *c.(*stubs.VCSConfig) // clone

			if s.ArbitraryAttribute != "" {
				// Note, we concat to the existing config here (not replace) so
				// that tests can verify that the right configs are made
				// available to Merge().
				cfg.ArbitraryAttribute += " + " + s.ArbitraryAttribute
			}

			if s.FilesystemPath != "" {
				cfg.FilesystemPath = s.FilesystemPath
			}

			if err := ctx.NormalizePath(&cfg.FilesystemPath); err != nil {
				return nil, err
			}

			return &cfg, nil
		},
	}
}
