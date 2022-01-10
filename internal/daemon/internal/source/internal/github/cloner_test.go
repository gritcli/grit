package github_test

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/common/config"
	"github.com/gritcli/grit/internal/daemon/internal/source/internal/git"
	. "github.com/gritcli/grit/internal/daemon/internal/source/internal/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.NewCloner()", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		configure func(*config.GitHub)
		driver    *Driver
		logger    logging.DiscardLogger
	)

	BeforeEach(func() {
		configure = func(*config.GitHub) {}
	})

	AfterEach(func() {
		cancel()
	})

	When("unauthenticated", func() {
		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachUnauthenticated(configure)
		})

		It("returns a git.Cloner", func() {
			cloner, dir, err := driver.NewCloner(ctx, gritPublicTestRepo.ID, logger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&git.Cloner{
				Config:       driver.Config.Git,
				SSHEndpoint:  "git@github.com:gritcli/test-public.git",
				HTTPEndpoint: "https://github.com/gritcli/test-public.git",
				Logger: logging.Tee(
					logging.Demote(
						logging.Prefix(
							driver.Logger,
							"clone[446260684]: ",
						),
					),
					logger,
				),
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})

	When("authenticated", func() {
		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachAuthenticated(configure)
		})

		It("returns a git.Cloner with the token as the HTTP password", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// Test with a private repository instead.
			cloner, dir, err := driver.NewCloner(ctx, gritPublicTestRepo.ID, logger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&git.Cloner{
				Config:       driver.Config.Git,
				SSHEndpoint:  "git@github.com:gritcli/test-public.git",
				HTTPEndpoint: "https://github.com/gritcli/test-public.git",
				HTTPPassword: driver.Config.Token,
				Logger: logging.Tee(
					logging.Demote(
						logging.Prefix(
							driver.Logger,
							"clone[446260684]: ",
						),
					),
					logger,
				),
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})
})
