package githubsource_test

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	publicOrgRepo = sourcedriver.RemoteRepo{
		ID:          "451303002",
		Name:        "grit-integration-tests-org/test-public",
		Description: "Used to test that Grit works with public GitHub repositories that belong to an organization.",
		WebURL:      "https://github.com/grit-integration-tests-org/test-public",
	}

	privateOrgRepo = sourcedriver.RemoteRepo{
		ID:          "451303236",
		Name:        "grit-integration-tests-org/test-private",
		Description: "Used to test that Grit works with private GitHub repositories that belong to an organization.",
		WebURL:      "https://github.com/grit-integration-tests-org/test-private",
	}

	publicUserRepo = sourcedriver.RemoteRepo{
		ID:          "451288349",
		Name:        "grit-integration-tests/test-public",
		Description: "Used to test that Grit works with public GitHub repositories.",
		WebURL:      "https://github.com/grit-integration-tests/test-public",
	}

	privateUserRepo = sourcedriver.RemoteRepo{
		ID:          "451288389",
		Name:        "grit-integration-tests/test-private",
		Description: "Used to test that Grit works with private GitHub repositories.",
		WebURL:      "https://github.com/grit-integration-tests/test-private",
	}

	// thirdPartyRepo is a repository that the authenticated user does not have
	// access to.
	//
	// The CI process uses a GitHub personal access token belonging to
	// @grit-integration-tests which is NOT a member of the "grit-cli"
	// organization.
	thirdPartyRepo = sourcedriver.RemoteRepo{
		ID:          "397822937",
		Name:        "gritcli/grit",
		Description: "Manage your local Git clones.",
		WebURL:      "https://github.com/gritcli/grit",
	}
)

var _ = Describe("func impl.Resolve()", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		drv    sourcedriver.Driver
		logger logging.DiscardLogger
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			ctx, cancel, drv = beforeEachUnauthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("does not resolve unqualified names", func() {
			repos, err := drv.Resolve(ctx, "test-public", logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("resolves an exact match using the API", func() {
			repos, err := drv.Resolve(ctx, publicUserRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(publicUserRepo))

			repos, err = drv.Resolve(ctx, publicOrgRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(publicOrgRepo))
		})

		It("returns nothing for a qualified name that does not exist", func() {
			repos, err := drv.Resolve(ctx, "grit-integration-tests/test-non-existant", logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("returns nothing for a qualified name that refers to a private repo", func() {
			repos, err := drv.Resolve(ctx, privateUserRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())

			repos, err = drv.Resolve(ctx, privateOrgRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			ctx, cancel, drv, _ = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("ignores invalid names", func() {
			repos, err := drv.Resolve(ctx, "has a space", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(BeEmpty())

			repos, err = drv.Resolve(ctx, "owner has a space/repo", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(BeEmpty())
		})

		It("resolves unqualified repo names using the cache", func() {
			repos, err := drv.Resolve(ctx, "test-public", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(publicUserRepo, publicOrgRepo))

			repos, err = drv.Resolve(ctx, "test-private", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(privateUserRepo, privateOrgRepo))
		})

		It("resolves an exact match using the cache", func() {
			repos, err := drv.Resolve(ctx, publicUserRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(publicUserRepo))

			repos, err = drv.Resolve(ctx, publicOrgRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(publicOrgRepo))
		})

		It("resolves an exact match for a private repo using the cache", func() {
			repos, err := drv.Resolve(ctx, privateUserRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(privateUserRepo))

			repos, err = drv.Resolve(ctx, privateOrgRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(privateOrgRepo))
		})

		It("resolves an exact match using the API", func() {
			repos, err := drv.Resolve(ctx, thirdPartyRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(thirdPartyRepo))
		})
	})
})
