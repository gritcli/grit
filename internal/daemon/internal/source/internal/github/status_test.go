package github_test

import (
	"context"

	"github.com/gritcli/grit/internal/daemon/internal/config"
	"github.com/gritcli/grit/plugin/driver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Driver", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		configure func(*config.GitHub)
		drv       driver.Driver
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			configure = func(*config.GitHub) {}
		})

		JustBeforeEach(func() {
			ctx, cancel, drv = beforeEachUnauthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		Describe("func Status()", func() {
			It("indicates that the user is unauthenticated", func() {
				status, err := drv.Status(ctx)
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
					status, err := drv.Status(ctx)
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
			ctx, cancel, drv = beforeEachAuthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		Describe("func Status()", func() {
			It("indicates that the user is authenticated", func() {
				status, err := drv.Status(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status).To(MatchRegexp(`@jmalloc, \d+ API requests remaining \(resets .+ from now\)`))
			})
		})
	})
})
