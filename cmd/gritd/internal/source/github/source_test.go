package github_test

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	. "github.com/gritcli/grit/cmd/gritd/internal/source/github"
	"github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Source", func() {
	var source source.Source

	BeforeEach(func() {
		var err error
		source, err = NewSource(
			"<source>",
			config.GitHubConfig{
				Domain: "github.example.com",
			},
			logging.SilentLogger,
		)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("func Name()", func() {
		It("returns the source name", func() {
			Expect(source.Name()).To(Equal("<source>"))
		})
	})

	Describe("func Description()", func() {
		It("returns a suitable description", func() {
			Expect(source.Description()).To(Equal("github.example.com (github enterprise)"))
		})
	})
})
