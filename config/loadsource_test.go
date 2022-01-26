package config_test

import (
	"errors"
	"reflect"

	. "github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/gritcli/grit/internal/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load() (source configuration)", func() {
	type aDifferentVCSConfigConcreteType struct {
		*stubs.VCSConfig
	}

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
			"driver-specific configuration",
			[]string{
				`source "test_source" "test_source_driver" {
					arbitrary_attribute = "<explicit>"
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "test_source",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/test_source",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<explicit>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default>",
						},
					},
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
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<implicit>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default>",
						},
					},
				},
			}),
			func(reg *registry.Registry) {
				loader := newSourceLoader()
				loader.ImplicitSourcesFunc = func(
					ctx sourcedriver.ConfigContext,
				) ([]sourcedriver.ImplicitSource, error) {
					cfg := &stubs.SourceConfig{
						ArbitraryAttribute: "<implicit>",
					}

					vcsConfig := &stubs.VCSConfig{}
					if err := ctx.UnmarshalVCSConfig(testVCSDriverName, &vcsConfig); err != nil {
						return nil, err
					}

					cfg.VCSs = map[string]vcsdriver.Config{
						testVCSDriverName: vcsConfig,
					}

					return []sourcedriver.ImplicitSource{
						{
							Name:   "implicit",
							Config: cfg,
						},
					}, nil
				}

				reg.RegisterSourceDriver(
					"test_source_driver_with_implicit_source",
					sourcedriver.Registration{
						Name:         "test_source_driver",
						ConfigLoader: loader,
					},
				)
			},
		),
		Entry(
			`implicit source does not override explicit source`,
			[]string{
				`source "implicit" "test_source_driver" {
					arbitrary_attribute = "<explicit>"
				}`,
			},
			withSource(defaultConfig, Source{
				Name:    "implicit",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/implicit",
				},
				Driver: &stubs.SourceConfig{
					ArbitraryAttribute: "<explicit>",
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: &stubs.VCSConfig{
							ArbitraryAttribute: "<default>",
						},
					},
				},
			}),
			func(reg *registry.Registry) {
				loader := newSourceLoader()
				loader.ImplicitSourcesFunc = func(
					ctx sourcedriver.ConfigContext,
				) ([]sourcedriver.ImplicitSource, error) {
					return []sourcedriver.ImplicitSource{
						{
							Name: "implicit",
							Config: &stubs.SourceConfig{
								ArbitraryAttribute: "<implicit>",
							},
						},
					}, nil
				}

				reg.RegisterSourceDriver(
					"test_source_driver_with_implicit_source",
					sourcedriver.Registration{
						Name:         "test_source_driver",
						ConfigLoader: loader,
					},
				)
			},
		),
		Entry(
			`disambiguate VCS drivers with the same name`,
			[]string{},
			withSource(defaultConfig, Source{
				Name:    "implicit",
				Enabled: true,
				Clones: Clones{
					Dir: "~/grit/implicit",
				},
				Driver: &stubs.SourceConfig{
					VCSs: map[string]vcsdriver.Config{
						testVCSDriverName: aDifferentVCSConfigConcreteType{},
					},
				},
			}),
			func(reg *registry.Registry) {
				// Register another VCS driver with the same name but a
				// different alias.
				reg.RegisterVCSDriver(
					testVCSDriverName+"_alias",
					vcsdriver.Registration{
						Name: testVCSDriverName,
						ConfigLoader: &stubs.VCSConfigLoader{
							DefaultsFunc: func(
								vcsdriver.ConfigContext,
							) (vcsdriver.Config, error) {
								return aDifferentVCSConfigConcreteType{}, nil
							},
						},
					},
				)

				loader := newSourceLoader()
				loader.ImplicitSourcesFunc = func(
					ctx sourcedriver.ConfigContext,
				) ([]sourcedriver.ImplicitSource, error) {
					var target aDifferentVCSConfigConcreteType
					if err := ctx.UnmarshalVCSConfig(testVCSDriverName, &target); err != nil {
						return nil, err
					}

					return []sourcedriver.ImplicitSource{
						{
							Name: "implicit",
							Config: &stubs.SourceConfig{
								VCSs: map[string]vcsdriver.Config{
									testVCSDriverName: target,
								},
							},
						},
					}, nil

				}

				reg.RegisterSourceDriver(
					"test_source_driver_uses_ambiguous_vcs",
					sourcedriver.Registration{
						Name:         "test_source_driver_uses_ambiguous_vcs",
						ConfigLoader: loader,
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
			`error producing implicit sources`,
			[]string{},
			`the implicit sources provided by the 'test_source_driver_with_implicit_source' driver cannot be loaded: <error>`,
			func(reg *registry.Registry) {
				loader := newSourceLoader()
				loader.ImplicitSourcesFunc = func(
					ctx sourcedriver.ConfigContext,
				) ([]sourcedriver.ImplicitSource, error) {
					return nil, errors.New("<error>")
				}

				reg.RegisterSourceDriver(
					"test_source_driver_with_implicit_source",
					sourcedriver.Registration{
						Name:         "test_source_driver",
						ConfigLoader: loader,
					},
				)
			},
		),
		Entry(
			`unrecognized VCS driver dependency`,
			[]string{
				`source "test_source" "test_source_driver_uses_unrecognized_vcs" {}`,
			},
			`the implicit sources provided by the 'test_source_driver_uses_unrecognized_vcs' driver cannot be loaded: dependency on unrecognized version control system ('<unrecognized>')`,
			func(reg *registry.Registry) {
				loader := newSourceLoader()
				loader.ImplicitSourcesFunc = func(
					ctx sourcedriver.ConfigContext,
				) ([]sourcedriver.ImplicitSource, error) {
					target := &stubs.VCSConfig{}
					return nil, ctx.UnmarshalVCSConfig("<unrecognized>", &target)
				}

				reg.RegisterSourceDriver(
					"test_source_driver_uses_unrecognized_vcs",
					sourcedriver.Registration{
						Name:         "test_source_driver_uses_unrecognized_vcs",
						ConfigLoader: loader,
					},
				)
			},
		),
		Entry(
			`VCS driver config type mismatch`,
			[]string{
				`source "test_source" "test_source_driver_uses_unrecognized_vcs" {}`,
			},
			`the implicit sources provided by the 'test_source_driver_uses_unrecognized_vcs' driver cannot be loaded: depends on incompatible version control system ('test_vcs_driver'), none of the matching drivers ('test_vcs_driver') use the same configuration structure`,
			func(reg *registry.Registry) {
				loader := newSourceLoader()
				loader.ImplicitSourcesFunc = func(
					ctx sourcedriver.ConfigContext,
				) ([]sourcedriver.ImplicitSource, error) {
					target := aDifferentVCSConfigConcreteType{}
					return nil, ctx.UnmarshalVCSConfig(testVCSDriverName, &target)
				}

				reg.RegisterSourceDriver(
					"test_source_driver_uses_unrecognized_vcs",
					sourcedriver.Registration{
						Name:         "test_source_driver_uses_unrecognized_vcs",
						ConfigLoader: loader,
					},
				)
			},
		),
		Entry(
			`configuration for unused VCS driver`,
			[]string{
				`source "test_source" "test_source_driver" {
					vcs "unused_vcs_driver" {}
				}`,
			},
			`<dir>/config-0.hcl: the 'test_source' source has configuration for the 'unused_vcs_driver' version control system but the source driver ('test_source_driver') does not support that VCS`,
			func(reg *registry.Registry) {
				reg.RegisterVCSDriver(
					"unused_vcs_driver",
					vcsdriver.Registration{
						Name:         "unused_vcs_driver",
						ConfigLoader: newVCSLoader(),
					},
				)
			},
		),
	)

	When("the source driver uses UnmarshalVCSConfig() incorrectly", func() {
		DescribeTable(
			"it panics",
			func(target interface{}, expect string) {
				reg := &registry.Registry{}
				reg.RegisterSourceDriver(
					"test_source_driver",
					sourcedriver.Registration{
						Name: "test_source_driver",
						ConfigLoader: &stubs.SourceConfigLoader{
							ImplicitSourcesFunc: func(
								ctx sourcedriver.ConfigContext,
							) ([]sourcedriver.ImplicitSource, error) {
								return nil, ctx.UnmarshalVCSConfig(testVCSDriverName, target)
							},
						},
					},
				)

				Expect(func() {
					Load("./does-not-exist", reg)
				}).To(PanicWith(expect))
			},
			Entry(
				`nil`,
				nil,
				`v must be a pointer to a concrete implementation of the vcsdriver.Config interface, but it is nil`,
			),
			Entry(
				`not a pointer`,
				"<not a pointer>",
				`v must be a pointer to a concrete implementation of the vcsdriver.Config interface, but string is not a pointer`,
			),
			Entry(
				`not a VCS config`,
				&struct{}{},
				`v must be a pointer to a concrete implementation of the vcsdriver.Config interface, but struct {} does not implement that interface`,
			),
			Entry(
				`not concrete`,
				reflect.New(
					reflect.TypeOf(
						(*vcsdriver.Config)(nil),
					).Elem(),
				).Interface(),
				`v must be a pointer to a concrete implementation of the vcsdriver.Config interface, but vcsdriver.Config is not a concrete type (it's an interface)`,
			),
		)
	})
})
