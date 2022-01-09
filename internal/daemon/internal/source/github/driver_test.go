package github_test

import (
	"context"
	"os"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/common/config"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	. "github.com/gritcli/grit/internal/daemon/internal/source/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type driver", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		configure func(*config.GitHub)
		driver    source.Driver
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			configure = func(*config.GitHub) {}
		})

		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachUnauthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		Describe("func Status()", func() {
			It("indicates that the user is unauthenticated", func() {
				status, err := driver.Status(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status).To(MatchRegexp(`unauthenticated, \d+ API requests remaining \(resets .+ from now\)`))
			})
		})

		When("unauthenticated due to invalid token", func() {
			BeforeEach(func() {
				configure = func(cfg *config.GitHub) {
					cfg.Token = "<invalid-token>"
				}
			})

			Describe("func Status()", func() {
				It("indicates that the token is invalid", func() {
					status, err := driver.Status(ctx)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status).To(Equal(`unauthenticated (invalid token)`))
				})
			})
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			configure = func(*config.GitHub) {}
		})

		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachAuthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		Describe("func Status()", func() {
			It("indicates that the user is authenticated", func() {
				status, err := driver.Status(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status).To(MatchRegexp(`@jmalloc, \d+ API requests remaining \(resets .+ from now\)`))
			})
		})
	})
})

// beforeEachAuthenticated returns the context and driver used for running
// integration tests with an authenticated user.
func beforeEachAuthenticated(configure ...func(*config.GitHub)) (context.Context, context.CancelFunc, source.Driver) {
	return initDriver(
		func() config.GitHub {
			token := os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")

			if token == "" {
				Skip("set GRIT_INTEGRATION_TEST_GITHUB_TOKEN to enable tests that use the GitHub API as an authenticated user")
			}

			return config.GitHub{
				Domain: "github.com",
				Token:  token,
			}
		},
		configure,
	)
}

// beforeEachAuthenticated returns the context and driver used for running
// integration tests without an authenticated user.
func beforeEachUnauthenticated(configure ...func(*config.GitHub)) (context.Context, context.CancelFunc, source.Driver) {
	return initDriver(
		func() config.GitHub {
			return config.GitHub{
				Domain: "github.com",
			}
		},
		configure,
	)
}

// initDriver creates and initializes a driver.
//
// The configuration is built starting with the result of cfg(), and then
// calling each function in configure in order to mutate the config as desired.
//
// It is intended for use in the beforeEachXXX() helper functions.
func initDriver(cfg func() config.GitHub, configure []func(*config.GitHub)) (context.Context, context.CancelFunc, source.Driver) {
	if os.Getenv("GRIT_INTEGRATION_TEST_GITHUB") == "" {
		Skip("set GRIT_INTEGRATION_TEST_GITHUB to enable tests that use the GitHub API")
	}

	c := cfg()
	for _, fn := range configure {
		fn(&c)
	}

	d := &Driver{
		Config: c,
		Logger: logging.SilentLogger,
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := d.Init(ctx); err != nil {
		cancel()
		skipIfRateLimited(err)
	}

	done := make(chan struct{})

	go func() {
		defer GinkgoRecover()
		defer close(done)

		err := d.Run(ctx)
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
	}, d
}

// skipIfRateLimited asserts that err is nil, or skips the test if err is a
// GitHub rate limit error.
func skipIfRateLimited(err error) {
	if _, ok := err.(*github.RateLimitError); ok {
		Skip("GitHub rate limit exceeded")
	}

	Expect(err).ShouldNot(HaveOccurred())
}
