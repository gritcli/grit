package config_test

import (
	. "github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("func Load() (source configuration)", func() {
	DescribeTable(
		"it returns the expected configuration",
		testLoadSuccess,
		Entry(
			"explicitly enabled source",
			[]string{
				`source "test_source" "test_source_driver" {
					enabled = true
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: sourceConfigStub{
					Value:     "<default>",
					VCSConfig: vcsConfigStub{Value: "<default>"},
				},
			}),
		),
		Entry(
			"explicitly disabled source",
			[]string{
				`source "test_source" "test_source_driver" {
					enabled = false
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: false,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: sourceConfigStub{
					Value:     "<default>",
					VCSConfig: vcsConfigStub{Value: "<default>"},
				},
			}),
		),
		Entry(
			"driver-specific configuration",
			[]string{
				`source "test_source" "test_source_driver" {
					value = "<explicit>"
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: sourceConfigStub{
					Value:     "<explicit>",
					VCSConfig: vcsConfigStub{Value: "<default>"},
				},
			}),
		),
		Entry(
			`implicit source`,
			[]string{},
			withSource(defaultConfig, Source{
				Name:    "implicit",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/implicit",
				},
				Driver: sourceConfigStub{
					Value:     "<implicit>",
					VCSConfig: vcsConfigStub{Value: "<default>"},
				},
			}),
			func(reg *registry.Registry) {
				reg.RegisterSourceDriver(
					"test_source_driver_with_implicit_source",
					sourcedriver.Registration{
						Name: "test_source_driver",
						NewConfigSchema: func() sourcedriver.ConfigSchema {
							return &sourceConfigSchemaStub{}
						},
						ImplicitSources: map[string]func() sourcedriver.ConfigSchema{
							"implicit": func() sourcedriver.ConfigSchema {
								return &sourceConfigSchemaStub{
									Value: "<implicit>",
								}
							},
						},
					},
				)
			},
		),
		Entry(
			`implicit source does not override explicit source`,
			[]string{
				`source "implicit" "test_source_driver" {
					value = "<explicit>"
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "implicit",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/implicit",
				},
				Driver: sourceConfigStub{
					Value:     "<explicit>",
					VCSConfig: vcsConfigStub{Value: "<default>"},
				},
			}),
			func(reg *registry.Registry) {
				reg.RegisterSourceDriver(
					"test_source_driver_with_implicit_source",
					sourcedriver.Registration{
						Name: "test_source_driver",
						NewConfigSchema: func() sourcedriver.ConfigSchema {
							return &sourceConfigSchemaStub{}
						},
						ImplicitSources: map[string]func() sourcedriver.ConfigSchema{
							"implicit": func() sourcedriver.ConfigSchema {
								return &sourceConfigSchemaStub{
									Value: "<implicit>",
								}
							},
						},
					},
				)
			},
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		testLoadFailure,
		Entry(
			`empty source name`,
			[]string{
				`source "" "test_source_driver" {}`,
			},
			`<dir>/config-0.hcl: source configurations must provide a name`,
		),
		Entry(
			`invalid source name`,
			[]string{
				`source "<invalid>" "test_source_driver" {}`,
			},
			`<dir>/config-0.hcl: '<invalid>' is not a valid source name, valid characters are ASCII letters, numbers and underscore`,
		),
		Entry(
			`duplicate source names`,
			[]string{
				`source "test_source" "test_source_driver" {}`,
				`source "test_source" "test_source_driver" {}`,
			},
			`<dir>/config-1.hcl: the 'test_source' source conflicts with a source of the same name in <dir>/config-0.hcl (source names are case-insensitive)`,
		),
		Entry(
			`duplicate source names (case-insensitive)`,
			[]string{
				`source "test_source" "test_source_driver" {}`,
				`source "TEST_SOURCE" "test_source_driver" {}`,
			},
			`<dir>/config-1.hcl: the 'TEST_SOURCE' source conflicts with a source of the same name in <dir>/config-0.hcl (source names are case-insensitive)`,
		),
		Entry(
			`empty driver name`,
			[]string{
				`source "test_source" "" {}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source has an empty driver name`,
		),
		Entry(
			`unrecognized source driver name`,
			[]string{
				`source "test_source" "<unrecognized>" {}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source uses an unrecognized driver ('<unrecognized>'), the supported source drivers are 'test_source_driver'`,
		),
		Entry(
			`source with a well-structured, but invalid body`,
			[]string{
				`source "test_source" "test_source_driver" {
					unrecognized = true
				}`,
			},
			`<dir>/config-0.hcl:2,6-18: Unsupported argument; An argument named "unrecognized" is not expected here.`,
		),
		Entry(
			`error normalizing driver configuration`,
			[]string{
				`source "test_source" "test_source_driver" {
					filesystem_path = "~someuser/path/to/nowhere"
				}`,
			},
			`<dir>/config-0.hcl: the configuration for the 'test_source' source cannot be loaded: cannot expand user-specific home dir`,
		),
		Entry(
			`error normalizing implicit source's driver configuration`,
			[]string{},
			`the configuration for the implicit 'implicit' source (provided by the 'test_source_driver_with_implicit_source' driver) cannot be loaded: cannot expand user-specific home dir`,
			func(reg *registry.Registry) {
				reg.RegisterSourceDriver(
					"test_source_driver_with_implicit_source",
					sourcedriver.Registration{
						Name: "test_source_driver",
						NewConfigSchema: func() sourcedriver.ConfigSchema {
							return &sourceConfigSchemaStub{}
						},
						ImplicitSources: map[string]func() sourcedriver.ConfigSchema{
							"implicit": func() sourcedriver.ConfigSchema {
								return &sourceConfigSchemaStub{
									FilesystemPath: "~someuser/path/to/nowhere",
								}
							},
						},
					},
				)
			},
		),
	)

	When("a source driver is implemented incorrectly", func() {
		When("unmarshaling a VCS config", func() {
			It("panics if the target is not a pointer", func() {
				reg := &registry.Registry{}

				reg.RegisterSourceDriver(
					"test_source_driver",
					sourcedriver.Registration{
						Name: "test_source_driver",
						NewConfigSchema: func() sourcedriver.ConfigSchema {
							return &sourceConfigSchemaStub{}
						},
						ImplicitSources: map[string]func() sourcedriver.ConfigSchema{
							"implicit": func() sourcedriver.ConfigSchema {
								return &sourceConfigSchemaStub{
									FilesystemPath: "~someuser/path/to/nowhere",
								}
							},
						},
					},
				)
			})
		})
	})
})
