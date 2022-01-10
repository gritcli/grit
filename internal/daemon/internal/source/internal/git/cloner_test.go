package git_test

import (
	"context"
	"os"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	. "github.com/gritcli/grit/internal/daemon/internal/source/internal/git"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Cloner", func() {
	var (
		ctx     context.Context
		cancel  context.CancelFunc
		logger  logging.BufferedLogger
		cloner  *Cloner
		tempDir string
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		logger.Reset()

		cloner = &Cloner{
			SSHEndpoint:  "git@github.com:gritcli/test-public.git",
			HTTPEndpoint: "https://github.com/gritcli/test-public.git",
			Logger:       &logger,
		}

		var err error
		tempDir, err = os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		cancel()

		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("func Clone()", func() {
		// Note, in order to minimise the number of clones, these tests do not
		// cover the full range of configuration options available. These are
		// covered in the tests for the git.useHTTP() function.

		It("clones via SSH using the SSH agent", func() {
			if os.Getenv("SSH_AUTH_SOCK") == "" {
				Skip("SSH agent is not available")
			}

			err := cloner.Clone(ctx, tempDir)
			Expect(err).ShouldNot(HaveOccurred())
			expectRemoteURL(tempDir, cloner.SSHEndpoint)
		})

		It("clones via SSH using an explicit private key", func() {
			cloner.Config.SSHKeyFile = "./testdata/deploy-key-no-passphrase"

			err := cloner.Clone(ctx, tempDir)
			Expect(err).ShouldNot(HaveOccurred())
			expectRemoteURL(tempDir, cloner.SSHEndpoint)
		})

		It("clones via SSH using an explicit private key with a passphrase", func() {
			cloner.Config.SSHKeyFile = "./testdata/deploy-key-with-passphrase"
			cloner.Config.SSHKeyPassphrase = "passphrase"

			err := cloner.Clone(ctx, tempDir)
			Expect(err).ShouldNot(HaveOccurred())
			expectRemoteURL(tempDir, cloner.SSHEndpoint)
		})

		It("clones via HTTP without authentication", func() {
			cloner.Config.PreferHTTP = true

			err := cloner.Clone(ctx, tempDir)
			Expect(err).ShouldNot(HaveOccurred())
			expectRemoteURL(tempDir, cloner.HTTPEndpoint)
		})

		It("clones via HTTP with authentication", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// Test cloning a private repository instead.
			token := os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")
			if token == "" {
				Skip("GRIT_INTEGRATION_TEST_GITHUB_TOKEN is not set")
			}

			cloner.Config.PreferHTTP = true
			cloner.HTTPPassword = token // username ignored by github

			err := cloner.Clone(ctx, tempDir)
			Expect(err).ShouldNot(HaveOccurred())
			expectRemoteURL(tempDir, cloner.HTTPEndpoint)
		})
	})
})

func expectRemoteURL(dir, url string) {
	repo, err := git.PlainOpen(dir)
	Expect(err).ShouldNot(HaveOccurred())

	rem, err := repo.Remote("origin")
	Expect(err).ShouldNot(HaveOccurred())
	Expect(rem.Config().URLs).To(ConsistOf(url))
}
