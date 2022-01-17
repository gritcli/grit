package github_test

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver/gitvcs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func impl.NewCloner()", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		drv    sourcedriver.Driver
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			ctx, cancel, drv = beforeEachUnauthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("returns a gitvcs.Cloner", func() {
			cloner, dir, err := drv.NewCloner(ctx, gritPublicTestRepo.ID, logging.SilentLogger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:  "git@github.com:gritcli/test-public.git",
				HTTPEndpoint: "https://github.com/gritcli/test-public.git",
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})

	When("authenticated", func() {
		var token string

		BeforeEach(func() {
			ctx, cancel, drv, token = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("returns a gitvcs.Cloner with the token as the HTTP password", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// Test with a private repository instead.
			cloner, dir, err := drv.NewCloner(ctx, gritPublicTestRepo.ID, logging.SilentLogger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:  "git@github.com:gritcli/test-public.git",
				HTTPEndpoint: "https://github.com/gritcli/test-public.git",
				HTTPPassword: token,
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})
})
