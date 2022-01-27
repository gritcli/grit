package githubsource_test

import (
	"context"

	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.Suggest()", func() {
	var (
		cancel context.CancelFunc
		src    sourcedriver.Source
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			_, cancel, src = beforeEachUnauthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("returns an empty slice", func() {
			repos := src.Suggest("")
			Expect(repos).To(BeEmpty())
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			_, cancel, src, _ = beforeEachAuthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		It("returns repositories with names that begin with the given word", func() {
			By("matching everything")

			repos := src.Suggest("")
			Expect(repos).To(ConsistOf(allTestRepos))

			By("matching part of the owner name")

			// match both @grit-integration-tests and @grit-integration-tests-org
			repos = src.Suggest("grit-integration-")
			Expect(repos).To(ConsistOf(allTestRepos))

			By("matching part of the unqualified repo name")

			repos = src.Suggest("test-pu")
			Expect(repos).To(ConsistOf(publicOrgRepo, publicUserRepo))

			By("matching part of the fully-qualified repo name")

			repos = src.Suggest("grit-integration-tests/test-pu")
			Expect(repos).To(ConsistOf(publicUserRepo))
		})
	})
})
