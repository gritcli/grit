package githubsource_test

import (
	"context"

	. "github.com/gritcli/grit/builtins/githubsource"
	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.Status()", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		src    sourcedriver.Source
	)

	When("unauthenticated", func() {
		var configure func(*Config)

		BeforeEach(func() {
			configure = func(*Config) {}
		})

		JustBeforeEach(func() {
			ctx, cancel, src = beforeEachUnauthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		It("indicates that the user is unauthenticated", func() {
			status, err := src.Status(ctx)
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
				status, err := src.Status(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status).To(Equal(`unauthenticated (invalid token)`))
			})
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			ctx, cancel, src, _ = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("indicates that the user is authenticated", func() {
			status, err := src.Status(ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status).To(MatchRegexp(`@grit-integration-tests, \d+ API requests remaining \(resets .+ from now\)`))
		})
	})
})
