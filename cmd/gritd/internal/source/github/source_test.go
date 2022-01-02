package github_test

import (
	"context"
	"os"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	. "github.com/gritcli/grit/cmd/gritd/internal/source/github"
	"github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type source", func() {
	var (
		src    source.Source
		logger logging.DiscardLogger
	)

	When("the source has not been initialized", func() {
		BeforeEach(func() {
			var err error
			src, err = NewSource("github-source", config.GitHubConfig{Domain: "github.com"}, logger)
			Expect(err).ShouldNot(HaveOccurred())
		})

		Describe("func Name()", func() {
			It("returns the source name", func() {
				Expect(src.Name()).To(Equal("github-source"))
			})
		})

		Describe("func Description()", func() {
			It("returns the server's domain name", func() {
				Expect(src.Description()).To(Equal("github.com"))
			})

			When("using GitHub Enterprise server", func() {
				BeforeEach(func() {
					var err error
					src, err = NewSource("github-source", config.GitHubConfig{Domain: "code.example.com"}, logger)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("explicitly states that GitHub Enterprise is being used", func() {
					Expect(src.Description()).To(Equal("code.example.com (github enterprise)"))
				})
			})
		})
	})

	When("the source has been initialized", func() {
		var cancel context.CancelFunc

		When("unauthenticated due to invalid token", func() {
			var originalToken = os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")

			BeforeEach(func() {
				os.Setenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN", "<invalid-token>")
				_, cancel, src = beforeEachUnauthenticated()
			})

			AfterEach(func() {
				cancel()
				os.Setenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN", originalToken)
			})

			It("works as in unauthenticated mode", func() {
				Expect(src.Description()).To(Equal("github.com")) // no username
			})
		})

		When("authenticated", func() {
			BeforeEach(func() {
				_, cancel, src = beforeEachAuthenticated()
			})

			AfterEach(func() {
				cancel()
			})

			Describe("func Description()", func() {
				It("includes the user name", func() {
					Expect(src.Description()).To(Equal("github.com (@jmalloc)"))
				})
			})
		})
	})
})

// beforeEachAuthenticated returns the context and source used for running
// integration tests with an authenticated user.
func beforeEachAuthenticated() (context.Context, context.CancelFunc, source.Source) {
	return initSource(func() config.GitHubConfig {
		token := os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")

		if token == "" {
			Skip("set GRIT_INTEGRATION_TEST_GITHUB_TOKEN to enable tests that use the GitHub API as an authenticated user")
		}

		return config.GitHubConfig{
			Domain: "github.com",
			Token:  token,
		}
	})
}

// beforeEachAuthenticated returns the context and source used for running
// integration tests without an authenticated user.
func beforeEachUnauthenticated() (context.Context, context.CancelFunc, source.Source) {
	return initSource(func() config.GitHubConfig {
		return config.GitHubConfig{
			Domain: "github.com",
		}
	})
}

// initSource creates and initializes a source using the config returned by
// cfg(). It is intended for use in the beforeEachXXX() helper functions.
func initSource(cfg func() config.GitHubConfig) (context.Context, context.CancelFunc, source.Source) {
	if os.Getenv("GRIT_INTEGRATION_TEST_GITHUB") == "" {
		Skip("set GRIT_INTEGRATION_TEST_GITHUB to enable tests that use the GitHub API")
	}

	src, err := NewSource("github", cfg(), logging.SilentLogger)
	Expect(err).ShouldNot(HaveOccurred())

	ctx, cancel := context.WithCancel(context.Background())

	err = src.Init(ctx)
	if err != nil {
		cancel()
		skipIfRateLimited(err)
	}

	done := make(chan struct{})

	go func() {
		defer GinkgoRecover()
		defer close(done)

		err := src.Run(ctx)
		if err != context.Canceled {
			Expect(err).ShouldNot(HaveOccurred())
		}
	}()

	return ctx, func() {
		cancel()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			Fail("timed out waiting for Run() goroutine to finish")
		}
	}, src
}

// skipIfRateLimited asserts that err is nil, or skips the test if err is a
// GitHub rate limit error.
func skipIfRateLimited(err error) {
	if _, ok := err.(*github.RateLimitError); ok {
		Skip("GitHub rate limit exceeded")
	}

	Expect(err).ShouldNot(HaveOccurred())
}
