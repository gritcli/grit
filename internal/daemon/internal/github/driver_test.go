package github_test

import (
	"context"
	"os"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	. "github.com/gritcli/grit/internal/daemon/internal/github"
	"github.com/gritcli/grit/plugin/sourcedriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// beforeEachAuthenticated returns the context and driver used for running
// integration tests with an authenticated user.
func beforeEachAuthenticated(configure ...func(*Config)) (
	_ context.Context,
	_ context.CancelFunc,
	_ sourcedriver.Driver,
	token string,
) {
	token = os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")

	ctx, cancel, drv := initDriver(
		func() Config {
			if token == "" {
				Skip("set GRIT_INTEGRATION_TEST_GITHUB_TOKEN to enable tests that use the GitHub API as an authenticated user")
			}

			return Config{
				Domain: "github.com",
				Token:  token,
			}
		},
		configure,
	)

	return ctx, cancel, drv, token
}

// beforeEachUnauthenticated returns the context and driver used for running
// integration tests without an authenticated user.
func beforeEachUnauthenticated(configure ...func(*Config)) (
	context.Context,
	context.CancelFunc,
	sourcedriver.Driver,
) {
	return initDriver(
		func() Config {
			return Config{
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
func initDriver(
	cfg func() Config,
	configure []func(*Config),
) (
	context.Context,
	context.CancelFunc,
	sourcedriver.Driver,
) {
	if os.Getenv("GRIT_INTEGRATION_TEST_USE_GITHUB_API") == "" {
		Skip("set GRIT_INTEGRATION_TEST_USE_GITHUB_API to enable tests that use the GitHub API")
	}

	c := cfg()
	for _, fn := range configure {
		fn(&c)
	}

	d := c.NewDriver()

	ctx, cancel := context.WithCancel(context.Background())

	if err := d.Init(ctx, logging.SilentLogger); err != nil {
		cancel()
		skipIfRateLimited(err)
	}

	done := make(chan struct{})

	go func() {
		defer GinkgoRecover()
		defer close(done)

		err := d.Run(ctx, logging.SilentLogger)
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
