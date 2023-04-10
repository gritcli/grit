package githubsource_test

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/v50/github"
	. "github.com/gritcli/grit/daemon/internal/builtins/githubsource"
	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// beforeEachAuthenticated returns the context and source used for running
// integration tests with an authenticated user.
func beforeEachAuthenticated(configure ...func(*Config)) (
	_ context.Context,
	_ context.CancelFunc,
	_ sourcedriver.Source,
	token string,
) {
	token = os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")

	ctx, cancel, s := initSource(
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

	return ctx, cancel, s, token
}

// beforeEachUnauthenticated returns the context and source used for running
// integration tests without an authenticated user.
func beforeEachUnauthenticated(configure ...func(*Config)) (
	context.Context,
	context.CancelFunc,
	sourcedriver.Source,
) {
	return initSource(
		func() Config {
			return Config{
				Domain: "github.com",
			}
		},
		configure,
	)
}

// initSource creates and initializes a source.
//
// The configuration is built starting with the result of cfg(), and then
// calling each function in configure in order to mutate the config as desired.
//
// It is intended for use in the beforeEachXXX() helper functions.
func initSource(
	cfg func() Config,
	configure []func(*Config),
) (
	context.Context,
	context.CancelFunc,
	sourcedriver.Source,
) {
	if os.Getenv("GRIT_INTEGRATION_TEST_USE_GITHUB_API") == "" {
		Skip("set GRIT_INTEGRATION_TEST_USE_GITHUB_API to enable tests that use the GitHub API")
	}

	c := cfg()
	for _, fn := range configure {
		fn(&c)
	}

	src := c.NewSource()

	ctx, cancel := context.WithCancel(context.Background())

	if err := src.Init(
		ctx,
		sourcedriver.InitParameters{},
		logs.Discard,
	); err != nil {
		cancel()
		skipIfRateLimited(err)
	}

	done := make(chan struct{})

	go func() {
		defer GinkgoRecover()
		defer close(done)

		err := src.Run(ctx, logs.Discard)
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
