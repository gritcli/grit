package githubsource_test

import (
	"context"

	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.Resolve()", func() {
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

		It("does not resolve unqualified names", func() {
			repos, err := src.Resolve(ctx, "test-public", logs.Discard)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("resolves an exact match using the API", func() {
			repos, err := src.Resolve(ctx, publicUserRepo.Name, logs.Discard)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(publicUserRepo))

			repos, err = src.Resolve(ctx, publicOrgRepo.Name, logs.Discard)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(publicOrgRepo))
		})

		It("returns nothing for a qualified name that does not exist", func() {
			repos, err := src.Resolve(ctx, "grit-integration-tests/test-non-existant", logs.Discard)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("returns nothing for a qualified name that refers to a private repo", func() {
			repos, err := src.Resolve(ctx, privateUserRepo.Name, logs.Discard)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())

			repos, err = src.Resolve(ctx, privateOrgRepo.Name, logs.Discard)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			ctx, cancel, src, _ = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("ignores invalid names", func() {
			repos, err := src.Resolve(ctx, "has a space", logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(BeEmpty())

			repos, err = src.Resolve(ctx, "owner has a space/repo", logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(BeEmpty())
		})

		It("resolves unqualified repo names using the cache", func() {
			repos, err := src.Resolve(ctx, "test-public", logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(publicUserRepo, publicOrgRepo))

			repos, err = src.Resolve(ctx, "test-private", logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(privateUserRepo, privateOrgRepo))
		})

		It("resolves an exact match using the cache", func() {
			repos, err := src.Resolve(ctx, publicUserRepo.Name, logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(publicUserRepo))

			repos, err = src.Resolve(ctx, publicOrgRepo.Name, logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(publicOrgRepo))
		})

		It("resolves an exact match for a private repo using the cache", func() {
			repos, err := src.Resolve(ctx, privateUserRepo.Name, logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(privateUserRepo))

			repos, err = src.Resolve(ctx, privateOrgRepo.Name, logs.Discard)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(privateOrgRepo))
		})

		It("resolves an exact match using the API", func() {
			repos, err := src.Resolve(ctx, thirdPartyRepo.Name, logs.Discard)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(thirdPartyRepo))
		})
	})
})
