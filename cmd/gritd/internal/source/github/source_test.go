package github_test

import (
	"context"
	"os"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	. "github.com/gritcli/grit/cmd/gritd/internal/source/github"
	"github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Source", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		source source.Source
		cfg    config.GitHubConfig
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		cfg = config.GitHubConfig{
			Domain: "github.com",
			Token:  os.Getenv("GRIT_TEST_GITHUB_TOKEN"),
		}
	})

	JustBeforeEach(func() {
		var err error
		source, err = NewSource(
			"github-source",
			cfg,
			logging.SilentLogger,
		)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		cancel()
	})

	Describe("func Name()", func() {
		It("returns the source name", func() {
			Expect(source.Name()).To(Equal("github-source"))
		})
	})

	When("the source has not been initialized", func() {
		Describe("func Description()", func() {
			It("returns the server's domain name", func() {
				Expect(source.Description()).To(Equal("github.com"))
			})

			When("using GitHub Enterprise server", func() {
				BeforeEach(func() {
					cfg.Domain = "code.example.com"
				})

				It("explicitly states that GitHub Enterprise is being used", func() {
					Expect(source.Description()).To(Equal("code.example.com (github enterprise)"))
				})
			})
		})
	})

	When("the source has been initialized", func() {
		JustBeforeEach(func() {
			err := source.Init(ctx)
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("authenticated", func() {
			BeforeEach(func() {
				if cfg.Token == "" {
					Skip("authentication token not available")
				}
			})

			Describe("func Description()", func() {
				It("includes the user name", func() {
					Expect(source.Description()).To(Equal("github.com (@jmalloc)"))
				})
			})
		})

		When("unauthenticated", func() {
			BeforeEach(func() {
				cfg.Token = ""
			})
		})
	})
})
