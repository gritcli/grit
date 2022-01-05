package github_test

import (
	"context"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dogmatiq/dodeca/logging"
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

		It("clones the repository", func() {
			err := src.Clone(ctx, gritRepo.ID, dir, logger)
			skipIfRateLimited(err)

			// We check that the clone was successful by checking that the
			// license file was obtained. We don't want to be too strict, as
			// there's no guarantee this test is running from the main branch.
			license, err := ioutil.ReadFile(filepath.Join(dir, "LICENSE"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(license)).To(ContainSubstring("MIT License"))
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
