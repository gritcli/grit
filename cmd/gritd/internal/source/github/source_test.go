package github_test

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
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
		src    source.Source
		cfg    config.GitHubConfig
		out    strings.Builder
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		cfg = config.GitHubConfig{
			Domain: "github.com",
			Token:  os.Getenv("GRIT_TEST_GITHUB_TOKEN"),
		}

		out.Reset()
	})

	JustBeforeEach(func() {
		var err error
		src, err = NewSource(
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
			Expect(src.Name()).To(Equal("github-source"))
		})
	})

	When("the source has not been initialized", func() {
		Describe("func Description()", func() {
			It("returns the server's domain name", func() {
				Expect(src.Description()).To(Equal("github.com"))
			})

			When("using GitHub Enterprise server", func() {
				BeforeEach(func() {
					cfg.Domain = "code.example.com"
				})

				It("explicitly states that GitHub Enterprise is being used", func() {
					Expect(src.Description()).To(Equal("code.example.com (github enterprise)"))
				})
			})
		})
	})

	When("the source has been initialized", func() {
		JustBeforeEach(func() {
			err := src.Init(ctx)
			skipIfRateLimited(err)
		})

		When("unauthenticated (invalid token)", func() {
			BeforeEach(func() {
				cfg.Token = "<invalid>"
			})

			It("works in unauthenticated mode", func() {
				Expect(src.Description()).To(Equal("github.com"))
			})
		})

		When("unauthenticated (no token)", func() {
			BeforeEach(func() {
				cfg.Token = ""
			})

			Describe("func Resolve()", func() {
				It("does not resolve unqualified names", func() {
					repos, err := src.Resolve(ctx, "grit", &out)
					skipIfRateLimited(err)
					Expect(repos).To(BeEmpty())
				})

				It("resolves an exact match using the API", func() {
					repos, err := src.Resolve(ctx, "gritcli/grit", &out)
					skipIfRateLimited(err)
					Expect(repos).To(ConsistOf(
						source.Repo{
							ID:          "397822937",
							Name:        "gritcli/grit",
							Description: "Manage your local Git clones.",
							WebURL:      "https://github.com/gritcli/grit",
						},
					))
				})

				It("returns nothing for a qualified name that does not exist", func() {
					repos, err := src.Resolve(ctx, "gritcli/non-existant", &out)
					skipIfRateLimited(err)
					Expect(repos).To(BeEmpty())
				})
			})
		})

		When("authenticated", func() {
			BeforeEach(func() {
				if cfg.Token == "" {
					Skip("authentication token not available")
				}
			})

			Describe("func Description()", func() {
				It("includes the user name", func() {
					Expect(src.Description()).To(Equal("github.com (@jmalloc)"))
				})
			})

			Describe("func Resolve()", func() {
				It("ignores invalid names", func() {
					repos, err := src.Resolve(ctx, "has a space", &out)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(repos).To(BeEmpty())
				})

				It("resolves unqualified repo names using the cache", func() {
					repos, err := src.Resolve(ctx, "grit", &out)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(repos).To(ConsistOf(
						source.Repo{
							ID:          "85247932",
							Name:        "jmalloc/grit",
							Description: "Keep track of your local Git clones.",
							WebURL:      "https://github.com/jmalloc/grit",
						},
						source.Repo{
							ID:          "397822937",
							Name:        "gritcli/grit",
							Description: "Manage your local Git clones.",
							WebURL:      "https://github.com/gritcli/grit",
						},
					))
				})

				It("resolves an exact match using the cache", func() {
					repos, err := src.Resolve(ctx, "gritcli/grit", &out)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(repos).To(ConsistOf(
						source.Repo{
							ID:          "397822937",
							Name:        "gritcli/grit",
							Description: "Manage your local Git clones.",
							WebURL:      "https://github.com/gritcli/grit",
						},
					))
				})

				It("resolves an exact match using the API", func() {
					// google/go-github this will never be in the cache for
					// @jmalloc (who owns the token used under CI)
					repos, err := src.Resolve(ctx, "google/go-github", &out)
					skipIfRateLimited(err)
					Expect(repos).To(ConsistOf(
						source.Repo{
							ID:          "10270722",
							Name:        "google/go-github",
							Description: "Go library for accessing the GitHub API",
							WebURL:      "https://github.com/google/go-github"},
					))
				})
			})
		})
	})
})

// skipIfRateLimited asserts that err is nil, or skips the test if err is a
// GitHub rate limit error.
func skipIfRateLimited(err error) {
	if _, ok := err.(*github.RateLimitError); ok {
		Skip("GitHub rate limit exceeded")
	}

	Expect(err).ShouldNot(HaveOccurred())
}
