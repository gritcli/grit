package githubsource_test

import (
	"context"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.Suggest()", func() {
	var (
		src sourcedriver.Source
	)

	When("unauthenticated", func() {
		BeforeEach(func() {
			var cancel context.CancelFunc
			_, cancel, src = beforeEachUnauthenticated()
			DeferCleanup(cancel)
		})

		It("returns an empty slice", func() {
			repos := src.Suggest("", logs.Discard)
			Expect(repos).To(BeEmpty())
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			var cancel context.CancelFunc
			_, cancel, src, _ = beforeEachAuthenticated()
			DeferCleanup(cancel)
		})

		It("returns repositories with names that begin with the given word", func() {
			By("matching everything")

			repos := src.Suggest("", logs.Discard)
			Expect(repos).To(Equal(
				map[string][]sourcedriver.RemoteRepo{
					publicOrgRepo.Name:   {publicOrgRepo},
					privateOrgRepo.Name:  {privateOrgRepo},
					publicUserRepo.Name:  {publicUserRepo},
					privateUserRepo.Name: {privateUserRepo},
				},
			))

			By("matching part of the owner name")

			// match both @grit-integration-tests and @grit-integration-tests-org
			repos = src.Suggest("grit-integration-", logs.Discard)
			Expect(repos).To(Equal(
				map[string][]sourcedriver.RemoteRepo{
					publicOrgRepo.Name:   {publicOrgRepo},
					privateOrgRepo.Name:  {privateOrgRepo},
					publicUserRepo.Name:  {publicUserRepo},
					privateUserRepo.Name: {privateUserRepo},
				},
			))

			By("matching part of the unqualified repo name")

			repos = src.Suggest("test-pu", logs.Discard)
			Expect(repos).To(
				HaveKeyWithValue(
					"test-public",
					ConsistOf(
						publicUserRepo,
						publicOrgRepo,
					),
				),
			)
			Expect(repos).To(HaveLen(1))

			By("matching part of the fully-qualified repo name")

			repos = src.Suggest("grit-integration-tests/test-pu", logs.Discard)
			Expect(repos).To(Equal(
				map[string][]sourcedriver.RemoteRepo{
					publicUserRepo.Name: {publicUserRepo},
				},
			))
		})
	})
})
