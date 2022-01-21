package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/gritcli/grit/internal/stubs"
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
	hooks ...func(r *registry.Registry),
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
	hooks ...func(r *registry.Registry),
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
func newRegistry() *registry.Registry {
	reg := &registry.Registry{}

	reg.RegisterSourceDriver(
		testSourceDriverName,
		sourcedriver.Registration{
			Name: testSourceDriverName,
			NewConfigSchema: func() sourcedriver.ConfigSchema {
				return newSourceStub()
			},
		},
	)

	reg.RegisterVCSDriver(
		testVCSDriverName,
		vcsdriver.Registration{
			Name: testVCSDriverName,
			NewConfigSchema: func() vcsdriver.ConfigSchema {
				return newVCSStub()
			},
		},
	)

	return reg
}

// newSourceStub returns a new stub of sourcedriver.ConfigSchema for testing
// source driver configuration.
func newSourceStub() *stubs.SourceDriverConfigSchema {
	return &stubs.SourceDriverConfigSchema{
		NormalizeFunc: func(
			nc sourcedriver.ConfigNormalizeContext,
			s *stubs.SourceDriverConfigSchema,
		) (sourcedriver.Config, error) {
			cfg := &stubs.SourceDriverConfig{
				ArbitraryAttribute: s.ArbitraryAttribute,
				FilesystemPath:     s.FilesystemPath,
			}

			if cfg.ArbitraryAttribute == "" {
				cfg.ArbitraryAttribute = "<default>"
			}

			if err := nc.NormalizePath(&cfg.FilesystemPath); err != nil {
				return nil, err
			}

			vcsConfig := &stubs.VCSDriverConfig{}
			if err := nc.UnmarshalVCSConfig(testVCSDriverName, &vcsConfig); err != nil {
				return nil, err
			}

			cfg.VCSs = map[string]vcsdriver.Config{
				testVCSDriverName: vcsConfig,
			}

			return cfg, nil
		},
	}
}

// newVCSStub returns a new stub of vcsdriver.ConfigSchema for testing VCS
// driver configuration.
func newVCSStub() *stubs.VCSDriverConfigSchema {
	return &stubs.VCSDriverConfigSchema{
		NormalizeGlobalsFunc: func(
			nc vcsdriver.ConfigNormalizeContext,
			s *stubs.VCSDriverConfigSchema,
		) (vcsdriver.Config, error) {
			cfg := &stubs.VCSDriverConfig{
				ArbitraryAttribute: s.ArbitraryAttribute,
				FilesystemPath:     s.FilesystemPath,
			}

			if cfg.ArbitraryAttribute == "" {
				cfg.ArbitraryAttribute = "<default>"
			}

			if err := nc.NormalizePath(&cfg.FilesystemPath); err != nil {
				return nil, err
			}

			return cfg, nil
		},

		NormalizeSourceSpecificFunc: func(
			nc vcsdriver.ConfigNormalizeContext,
			g vcsdriver.Config,
			s *stubs.VCSDriverConfigSchema,
		) (vcsdriver.Config, error) {
			cfg := *g.(*stubs.VCSDriverConfig) // clone

			if s.ArbitraryAttribute != "" {
				// Note, we concat to the default here (not replace) so that
				// tests can verify that the defaults are made available to
				// NormalizeSourceSpecific()
				cfg.ArbitraryAttribute += s.ArbitraryAttribute
			}

			if s.FilesystemPath != "" {
				cfg.FilesystemPath = s.FilesystemPath
			}

			if err := nc.NormalizePath(&cfg.FilesystemPath); err != nil {
				return nil, err
			}

			return &cfg, nil
		},
	}
}
