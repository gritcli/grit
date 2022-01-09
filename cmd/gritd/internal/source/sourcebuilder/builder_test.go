package sourcebuilder_test

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/cmd/gritd/internal/source/github"
	. "github.com/gritcli/grit/cmd/gritd/internal/source/sourcebuilder"
	"github.com/gritcli/grit/internal/config"
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
					Name: "test-source",
					Clones: config.Clones{
						Dir: "/path/to/clones",
					},
					DriverConfig: config.GitHub{
						Domain: "github.com",
					},
				},
				source.Source{
					Name:        "test-source",
					Description: "github.com",
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
