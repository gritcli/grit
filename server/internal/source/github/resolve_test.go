package github_test

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/server/internal/source"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	gritRepo = source.Repo{
		ID:          "397822937",
		Name:        "gritcli/grit",
		Description: "Manage your local Git clones.",
		WebURL:      "https://github.com/gritcli/grit",
	}

	gritPrivateTestRepo = source.Repo{
		ID:          "445039240",
		Name:        "gritcli/test-private",
		Description: "Used to test that Grit works with private GitHub repositories.",
		WebURL:      "https://github.com/gritcli/test-private",
	}

	gritV1Repo = source.Repo{
		ID:          "85247932",
		Name:        "jmalloc/grit",
		Description: "Keep track of your local Git clones.",
		WebURL:      "https://github.com/jmalloc/grit",
	}

	// thirdPartyRepo is a repo that the authenticated user does not have access
	// to. The CI process currently uses a GitHub personal access token
	// belonging to @jmalloc, who presumably would never be granted access to
	// anything in the "google" organization ;)
	thirdPartyRepo = source.Repo{
		ID:          "10270722",
		Name:        "google/go-github",
		Description: "Go library for accessing the GitHub API",
		WebURL:      "https://github.com/google/go-github",
	}
)

var _ = Describe("func source.Resolve()", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		driver source.Driver
		logger logging.DiscardLogger
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			ctx, cancel, driver = beforeEachUnauthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("does not resolve unqualified names", func() {
			repos, err := driver.Resolve(ctx, "grit", logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("resolves an exact match using the API", func() {
			repos, err := driver.Resolve(ctx, "gritcli/grit", logger)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(gritRepo))
		})

		It("returns nothing for a qualified name that does not exist", func() {
			repos, err := driver.Resolve(ctx, "gritcli/test-non-existant", logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})

		It("returns nothing for a qualified name that refers to a private repo", func() {
			repos, err := driver.Resolve(ctx, gritPrivateTestRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(BeEmpty())
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			ctx, cancel, driver = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("ignores invalid names", func() {
			repos, err := driver.Resolve(ctx, "has a space", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(BeEmpty())

			repos, err = driver.Resolve(ctx, "owner has a space/repo", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(BeEmpty())
		})

		It("resolves unqualified repo names using the cache", func() {
			repos, err := driver.Resolve(ctx, "grit", logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(gritRepo, gritV1Repo))
		})

		It("resolves an exact match using the cache", func() {
			repos, err := driver.Resolve(ctx, gritRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(gritRepo))
		})

		XIt("resolves an exact match for a private repo using the cache", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// This requires the "repo" scope on the personal-access-token which
			// grants read/write access to private repos. We currently use
			// @jmalloc's access token in GHA, so this is not feasible. We need
			// to setup a user specifically for testing this.
			repos, err := driver.Resolve(ctx, gritPrivateTestRepo.Name, logger)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repos).To(ConsistOf(gritPrivateTestRepo))
		})

		It("resolves an exact match using the API", func() {
			repos, err := driver.Resolve(ctx, thirdPartyRepo.Name, logger)
			skipIfRateLimited(err)
			Expect(repos).To(ConsistOf(thirdPartyRepo))
		})
	})
})
