package sourcebuilder_test

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/common/config"
	"github.com/gritcli/grit/server/internal/source"
	"github.com/gritcli/grit/server/internal/source/github"
	. "github.com/gritcli/grit/server/internal/source/sourcebuilder"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Listen()", func() {
	var (
		logger  logging.BufferedLogger
		builder *Builder
	)

	BeforeEach(func() {
		logger.Reset()

		builder = &Builder{
			Logger: &logger,
		}
	})

	Describe("func FromConfig()", func() {
		It("constructs all enabled sources from the configuration", func() {
			cfg := config.Config{
				Sources: map[string]config.Source{
					"github-test-source": {
						Name:    "github-test-source",
						Enabled: true,
						Clones: config.Clones{
							Dir: "/path/to/clones/github",
						},
						Driver: config.GitHub{
							Domain: "github.com",
						},
					},
					"disabled-test-source": {
						Name:    "disabled-test-source",
						Enabled: false,
						// None of the other fields are inspected at all if the
						// source is disabled.
					},
				},
			}

			sources := builder.FromConfig(cfg)
			Expect(sources).To(ConsistOf(
				source.Source{
					Name:        "github-test-source",
					Description: "github.com",
					CloneDir:    "/path/to/clones/github",
					Driver: &github.Driver{
						Config: config.GitHub{
							Domain: "github.com",
						},
						Logger: logging.Prefix(&logger, "source[github-test-source]: "),
					},
				},
			))
		})
	})

	Describe("func FromSourceConfig()", func() {
		table.DescribeTable(
			"it constructs sources based on the driver configuration type",
			func(cfg config.Source, expect source.Source) {
				src := builder.FromSourceConfig(cfg)
				Expect(src).To(Equal(expect))
			},
			Entry(
				"github",
				config.Source{
					Name:    "test-source",
					Enabled: false, // note, this is not checked by FromSourceConfig()
					Clones: config.Clones{
						Dir: "/path/to/clones",
					},
					Driver: config.GitHub{
						Domain: "github.com",
					},
				},
				source.Source{
					Name:        "test-source",
					Description: "github.com",
					CloneDir:    "/path/to/clones",
					Driver: &github.Driver{
						Config: config.GitHub{
							Domain: "github.com",
						},
						Logger: logging.Prefix(&logger, "source[test-source]: "),
					},
				},
			),
		)
	})
})