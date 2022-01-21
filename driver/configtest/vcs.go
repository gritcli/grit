package configtest

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// VCSDriverTest is a test that tests VCS driver configuration.
type VCSDriverTest = table.TableEntry

// TestVCSDriver runs a series of tests
func TestVCSDriver(
	r vcsdriver.Registration,
	zero vcsdriver.Config,
	tests ...VCSDriverTest,
) {
	table.DescribeTable(
		"it loads the configuration",
		func(
			content []string,
			expect func(dir string, cfg config.Config, err error),
		) {
			reg := &registry.Registry{}
			reg.RegisterVCSDriver(r.Name, r)

			reg.RegisterSourceDriver(
				"driver_under_test",
				sourcedriver.Registration{
					Name: "driver_under_test",
					NewConfigSchema: func() sourcedriver.ConfigSchema {
						return &vcsTestSourceConfigSchema{
							driverName: r.Name,
							unmarshalTarget: reflect.New(
								reflect.TypeOf(zero),
							).Interface(),
						}
					},
				},
			)

			dir, cleanup := writeConfigs(
				content...,
			)
			defer cleanup()

			cfg, err := config.Load(dir, reg)
			expect(dir, cfg, err)
		},
		tests...,
	)
}

// VCSSuccess returns a test that tests a default VCS driver configuration that
// is expected to pass.
func VCSSuccess(
	description string,
	defaultContent string,
	expect vcsdriver.Config,
) VCSDriverTest {
	return table.Entry(
		description,
		[]string{
			defaultContent,
			`source "test_source" "driver_under_test" {}`,
		},
		func(dir string, cfg config.Config, err error) {
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			for _, src := range cfg.Sources {
				if src.Name == "test_source" {
					cfg := src.Driver.(vcsTestSourceConfig).VCSConfig
					gomega.Expect(cfg).To(gomega.Equal(expect))
					return
				}
			}

			ginkgo.Fail("expected source was not defined")
		},
	)
}

// VCSSourceSpecificSuccess returns a test that tests a source-specific VCS
// driver configuration that is expected to pass.
func VCSSourceSpecificSuccess(
	description string,
	defaultContent string,
	sourceSpecificContent string,
	expect vcsdriver.Config,
) VCSDriverTest {
	return table.Entry(
		description,
		[]string{
			defaultContent,
			fmt.Sprintf(
				`source "test_source" "driver_under_test" {
				%s
				}`,
				sourceSpecificContent,
			),
		},
		func(dir string, cfg config.Config, err error) {
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			for _, src := range cfg.Sources {
				if src.Name == "test_source" {
					cfg := src.Driver.(vcsTestSourceConfig).VCSConfig
					gomega.Expect(cfg).To(gomega.Equal(expect))
					return
				}
			}

			ginkgo.Fail("expected source was not defined")
		},
	)
}

// VCSFailure returns a test that tests a VCS driver configuration that is
// expected to fail.
func VCSFailure(
	description string,
	content string,
	expect string,
) VCSDriverTest {
	return table.Entry(
		description,
		[]string{content},
		func(dir string, cfg config.Config, err error) {
			orig := format.TruncatedDiff
			format.TruncatedDiff = false
			defer func() {
				format.TruncatedDiff = orig
			}()

			message := strings.ReplaceAll(err.Error(), dir, "<dir>")
			gomega.Expect(message).To(gomega.Equal(expect))
		},
	)
}

type vcsTestSourceConfigSchema struct {
	driverName      string
	unmarshalTarget interface{}
}

func (s *vcsTestSourceConfigSchema) Normalize(
	nc sourcedriver.ConfigNormalizeContext,
) (sourcedriver.Config, error) {
	if err := nc.UnmarshalVCSConfig(
		s.driverName,
		s.unmarshalTarget,
	); err != nil {
		return nil, err
	}

	return vcsTestSourceConfig{
		VCSConfig: reflect.
			ValueOf(s.unmarshalTarget).
			Elem().
			Interface().(vcsdriver.Config),
	}, nil
}

type vcsTestSourceConfig struct {
	VCSConfig vcsdriver.Config
}

// NewDriver constructs a new driver that uses this configuration.
func (vcsTestSourceConfig) NewDriver() sourcedriver.Driver {
	panic("not implemented")
}

// DescribeSourceConfig returns a human-readable description of the
// configuration.
func (vcsTestSourceConfig) DescribeSourceConfig() string {
	panic("not implemented")
}
