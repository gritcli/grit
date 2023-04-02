package githubsource_test

import (
	"context"

	"github.com/gritcli/grit/builtins/gitvcs"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.NewCloner()", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		src    sourcedriver.Source
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			ctx, cancel, src = beforeEachUnauthenticated()
		})

		AfterEach(func() {
			cancel()
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
			ctx, cancel, src, token = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
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
