package github_test

import (
	"context"
	"math"
	"os"
	"strconv"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.Clone()", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		src    source.Source
		dir    string
		logger logging.DiscardLogger
	)

	BeforeEach(func() {
		var err error
		dir, err = os.MkdirTemp("", "grit-clone-test-")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		if dir != "" {
			os.RemoveAll(dir)
		}
	})

	When("unauthenticated", func() {
		BeforeEach(func() {
			ctx, cancel, src = beforeEachUnauthenticated()
		})

		AfterEach(func() {
			cancel()
		})

		When("the SSH agent is unavailable", func() {
			var orig string
			BeforeEach(func() {
				orig = os.Getenv("SSH_AUTH_SOCK")
				os.Setenv("SSH_AUTH_SOCK", "")
			})

			AfterEach(func() {
				os.Setenv("SSH_AUTH_SOCK", orig)
			})

			It("clones the repository using HTTP", func() {
				err := src.Clone(ctx, gritRepo.ID, dir, logger)
				skipIfRateLimited(err)

				repo, err := git.PlainOpen(dir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("https://github.com/gritcli/grit.git"))
			})
		})

		When("the SSH agent is available", func() {
			BeforeEach(func() {
				if os.Getenv("SSH_AUTH_SOCK") == "" {
					Skip("the SSH agent is not available")
				}

				// TODO: https://github.com/gritcli/grit/issues/12
				// Make this test work under CI.
			})

			It("clones the repository using SSH", func() {
				// This test asserts that even though _GRIT_ is not
				// authenticated with GitHub, the user's SSH key can still be
				// used for cloning.
				err := src.Clone(ctx, gritRepo.ID, dir, logger)
				skipIfRateLimited(err)

				repo, err := git.PlainOpen(dir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("git@github.com:gritcli/grit.git"))
			})
		})

		It("returns an error if the repository does not exist", func() {
			err := src.Clone(ctx, strconv.FormatInt(math.MaxInt64, 10), dir, logger)
			Expect(err).To(MatchError("GET https://api.github.com/repositories/9223372036854775807: 404 Not Found []"))
		})

		It("returns an error if the repository ID is invalid", func() {
			err := src.Clone(ctx, "<invalid>", dir, logger)
			Expect(err).To(MatchError("invalid repo ID, expected positive integer"))
		})

		It("returns an error if the repository ID is non-positive", func() {
			err := src.Clone(ctx, "0", dir, logger)
			Expect(err).To(MatchError("invalid repo ID, expected positive integer"))
		})
	})
})
