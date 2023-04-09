package githubsource_test

import (
	"context"

	"github.com/gritcli/grit/daemon/internal/builtins/gitvcs"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/logs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.NewCloner()", func() {
	var (
		ctx context.Context
		src sourcedriver.Source
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			var cancel context.CancelFunc
			ctx, cancel, src = beforeEachUnauthenticated()
			DeferCleanup(cancel)
		})

		It("returns a gitvcs.Cloner", func() {
			cloner, repo, err := src.NewCloner(ctx, publicUserRepo.ID, logs.Discard)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:  "git@github.com:grit-integration-tests/test-public.git",
				HTTPEndpoint: "https://github.com/grit-integration-tests/test-public.git",
			}))

			Expect(repo).To(Equal(publicUserRepo))
		})
	})

	When("authenticated", func() {
		var token string

		BeforeEach(func() {
			var cancel context.CancelFunc
			ctx, cancel, src, token = beforeEachAuthenticated()
			DeferCleanup(cancel)
		})

		It("returns a gitvcs.Cloner with the token as the HTTP password", func() {
			cloner, repo, err := src.NewCloner(ctx, privateUserRepo.ID, logs.Discard)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:  "git@github.com:grit-integration-tests/test-private.git",
				HTTPEndpoint: "https://github.com/grit-integration-tests/test-private.git",
				HTTPPassword: token,
			}))

			Expect(repo).To(Equal(privateUserRepo))
		})
	})
})
