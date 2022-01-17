package github_test

import (
	"context"

	. "github.com/gritcli/grit/internal/daemon/internal/github"
	"github.com/gritcli/grit/plugin/driver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func impl.Status()", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		drv    driver.Driver
	)

	When("unauthenticated", func() {
		var configure func(*Config)

		BeforeEach(func() {
			configure = func(*Config) {}
		})

		JustBeforeEach(func() {
			ctx, cancel, drv = beforeEachUnauthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		It("indicates that the user is unauthenticated", func() {
			status, err := drv.Status(ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status).To(MatchRegexp(`unauthenticated, \d+ API requests remaining \(resets .+ from now\)`))
		})

		When("unauthenticated due to invalid token", func() {
			BeforeEach(func() {
				configure = func(cfg *Config) {
					cfg.Token = "<invalid-token>"
				}
			})

			It("indicates that the token is invalid", func() {
				status, err := drv.Status(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status).To(Equal(`unauthenticated (invalid token)`))
			})
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			ctx, cancel, drv, _ = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("indicates that the user is authenticated", func() {
			status, err := drv.Status(ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status).To(MatchRegexp(`@jmalloc, \d+ API requests remaining \(resets .+ from now\)`))
		})
	})
})
