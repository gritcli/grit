package githubsource_test

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	gritRepo = sourcedriver.RemoteRepo{
		ID:          "397822937",
		Name:        "gritcli/grit",
		Description: "Manage your local Git clones.",
		WebURL:      "https://github.com/gritcli/grit",
	}

	gritPublicTestRepo = sourcedriver.RemoteRepo{
		ID:          "446260684",
		Name:        "gritcli/test-public",
		Description: "Used to test that Grit works with public GitHub repositories.",
		WebURL:      "https://github.com/gritcli/test-public",
	}

	gritPrivateTestRepo = sourcedriver.RemoteRepo{
		ID:          "445039240",
		Name:        "gritcli/test-private",
		Description: "Used to test that Grit works with private GitHub repositories.",
		WebURL:      "https://github.com/gritcli/test-private",
	}

	gritV1Repo = sourcedriver.RemoteRepo{
		ID:          "85247932",
		Name:        "jmalloc/grit",
		Description: "Keep track of your local Git clones.",
		WebURL:      "https://github.com/jmalloc/grit",
	}

	// thirdPartyRepo is a repo that the authenticated user does not have access
	// to. The CI process currently uses a GitHub personal access token
	// belonging to @jmalloc, who presumably would never be granted access to
	// anything in the "google" organization ;)
	thirdPartyRepo = sourcedriver.RemoteRepo{
		ID:          "10270722",
		Name:        "google/go-github",
		Description: "Go library for accessing the GitHub API",
		WebURL:      "https://github.com/google/go-github",
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
			repos, err := drv.Resolve(ctx, "grit", logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("resolves an exact match using the API", func() {
			repos, err := drv.Resolve(ctx, gritPublicTestRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(gritPublicTestRepo))
		})

		It("returns nothing for a qualified name that does not exist", func() {
			repos, err := drv.Resolve(ctx, "gritcli/test-non-existant", logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("returns nothing for a qualified name that refers to a private repo", func() {
			repos, err := drv.Resolve(ctx, gritPrivateTestRepo.Name, logger)
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
			repos, err := drv.Resolve(ctx, "grit", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(gritRepo, gritV1Repo))
		})

		It("resolves an exact match using the cache", func() {
			repos, err := drv.Resolve(ctx, gritPublicTestRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(gritPublicTestRepo))
		})

		XIt("resolves an exact match for a private repo using the cache", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// This requires the "repo" scope on the personal-access-token which
			// grants read/write access to private repos. We currently use
			// @jmalloc's access token in GHA, so this is not feasible. We need
			// to setup a user specifically for testing this.
			repos, err := drv.Resolve(ctx, gritPrivateTestRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(gritPrivateTestRepo))
		})

		It("resolves an exact match using the API", func() {
			repos, err := drv.Resolve(ctx, thirdPartyRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(thirdPartyRepo))
		})
	})
})
