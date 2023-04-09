package config_test

import (
	. "github.com/gritcli/grit/daemon/internal/config"
	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
	"github.com/gritcli/grit/daemon/internal/stubs"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("func Load() (clones configuration)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"sources use a directory within the default clone directory by default",
			[]string{
				`clones {
					dir = "/path/to/clones"
				}

				source "test_source" "test_source_driver" {}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "/path/to/clones/test_source",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<default>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default>",
						},
					},
				},
			}),
		),
		Entry(
			"sources can specify a clones configuration",
			[]string{
				`source "test_source" "test_source_driver" {
					clones {
						dir = "/path/to/clones"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "/path/to/clones",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<default>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default>",
						},
					},
				},
			}),
		),
		Entry(
			"sources can override the default clone directory",
			[]string{
				`clones {
					dir = "/path/to/clones"
				}

				source "test_source" "test_source_driver" {
					clones {
						dir = "/path/to/elsewhere"
					}
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "/path/to/elsewhere",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<default>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default>",
						},
					},
				},
			}),
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		testLoadFailure,
		Entry(
			`multiple files with clones defaults blocks`,
			[]string{
				`clones {}`,
				`clones {}`,
			},
			`<dir>/config-1.hcl: the global clones configuration is already defined in <dir>/config-0.hcl`,
		),
		Entry(
			`unexpandable global clones directory`,
			[]string{
				`clones {
					dir = "~someuser/path/to/clones"
				}`,
			},
			`<dir>/config-0.hcl: unable to resolve global clones directory: cannot expand user-specific home dir (~someuser/path/to/clones)`,
		),
		Entry(
			`unexpandable source-specific clones directory`,
			[]string{
				`source "test_source" "test_source_driver" {
					clones {
						dir = "~someuser/path/to/clones"
					}
				}`,
			},
			`<dir>/config-0.hcl: unable to resolve clones directory for the 'test_source' source: cannot expand user-specific home dir (~someuser/path/to/clones)`,
		),
	)

	Context("when the default global clones directory cannot be resolved", func() {
		var original string

		BeforeEach(func() {
			// HACK: We really shouldn't manipulate (or even have) global
			// variables like this, but it's the only cross-platform way to
			// force the home directory resolution to fail.
			original = DefaultClonesDirectory
			DefaultClonesDirectory = "~someuser/path/to/socket"
			DeferCleanup(func() {
				DefaultClonesDirectory = original
			})
		})

		DescribeTable(
			"it returns an error",
			testLoadFailure,
			Entry(
				`unexpandable default daemon socket`,
				[]string{},
				`unable to resolve default global clones directory: cannot expand user-specific home dir (~someuser/path/to/socket)`,
			),
		)
	})
})
