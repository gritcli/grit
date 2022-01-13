package gitvcs // note: no _test suffix to allow testing unexported useHTTP() method.

import (
	"context"
	"os"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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
		// covered in the tests for the useHTTP() method.

		It("clones via SSH using the SSH agent", func() {
			if os.Getenv("SSH_AUTH_SOCK") == "" {
				Skip("SSH agent is not available")
			}

			err := cloner.Clone(ctx, tempDir, &logger)
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.SSHEndpoint, &logger)
		})

		It("clones via SSH using an explicit private key", func() {
			cloner.SSHKeyFile = "./testdata/deploy-key-no-passphrase"

			err := cloner.Clone(ctx, tempDir, &logger)
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.SSHEndpoint, &logger)
		})

		It("clones via SSH using an explicit private key with a passphrase", func() {
			cloner.SSHKeyFile = "./testdata/deploy-key-with-passphrase"
			cloner.SSHKeyPassphrase = "passphrase"

			err := cloner.Clone(ctx, tempDir, &logger)
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.SSHEndpoint, &logger)
		})

		It("clones via HTTP without authentication", func() {
			cloner.PreferHTTP = true

			err := cloner.Clone(ctx, tempDir, &logger)
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.HTTPEndpoint, &logger)
		})

		It("clones via HTTP with authentication", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// Test cloning a private repository instead.
			token := os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")
			if token == "" {
				Skip("GRIT_INTEGRATION_TEST_GITHUB_TOKEN is not set")
			}

			cloner.PreferHTTP = true
			cloner.HTTPPassword = token // username ignored by github

			err := cloner.Clone(ctx, tempDir, &logger)
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.HTTPEndpoint, &logger)
		})
	})

	Describe("func useHTTP()", func() {
		DescribeTable(
			"it chooses the best available protocol",
			func(
				hasSSH, hasHTTP,
				hasSSHAgent, hasSSHPrivateKey, preferHTTP bool,
				expect string, // "ssh", "http" or an error message
			) {
				if !hasSSH {
					cloner.SSHEndpoint = ""
				}

				if !hasHTTP {
					cloner.HTTPEndpoint = ""
				}

				cloner.PreferHTTP = preferHTTP

				if hasSSHAgent {
					if os.Getenv("SSH_AUTH_SOCK") == "" {
						Skip("the SSH agent is unavailable")
					}
				} else {
					orig := os.Getenv("SSH_AUTH_SOCK")
					os.Setenv("SSH_AUTH_SOCK", "")
					defer os.Setenv("SSH_AUTH_SOCK", orig)
				}

				if hasSSHPrivateKey {
					cloner.SSHKeyFile = "<not empty>"
				}

				useHTTP, err := cloner.useHTTP()

				if expect != "ssh" && expect != "http" {
					Expect(err).To(MatchError(expect))
					return
				}

				Expect(err).ShouldNot(HaveOccurred())
				Expect(useHTTP).To(Equal(expect == "http"))
			},

			// Neither SSH nor HTTP URLs provided.
			Entry(
				`[ ] ssh [ ] http -- [ ] ssh agent [ ] ssh key [ ] prefer http`,
				false, false,
				false, false, false,
				"neither the SSH nor HTTP protocol is available",
			),

			// Only HTTP URL provided.
			Entry(
				`[ ] ssh [x] http -- [ ] ssh agent [ ] ssh key [ ] prefer http`,
				false, true,
				false, false, false,
				"http",
			),
			Entry(
				`[ ] ssh [x] http -- [ ] ssh agent [ ] ssh key [x] prefer http`,
				false, true,
				false, false, true,
				"http",
			),
			Entry(
				`[ ] ssh [x] http -- [ ] ssh agent [x] ssh key [ ] prefer http`,
				false, true,
				false, true, false,
				"http",
			),
			Entry(
				`[ ] ssh [x] http -- [ ] ssh agent [x] ssh key [x] prefer http`,
				false, true,
				false, true, true,
				"http",
			),
			Entry(
				`[ ] ssh [x] http -- [x] ssh agent [ ] ssh key [ ] prefer http`,
				false, true,
				true, false, false,
				"http",
			),
			Entry(
				`[ ] ssh [x] http -- [x] ssh agent [ ] ssh key [x] prefer http`,
				false, true,
				true, false, true,
				"http",
			),
			Entry(
				`[ ] ssh [x] http -- [x] ssh agent [x] ssh key [ ] prefer http`,
				false, true,
				true, true, false,
				"http",
			),
			Entry(
				`[ ] ssh [x] http -- [x] ssh agent [x] ssh key [x] prefer http`,
				false, true,
				true, true, true,
				"http",
			),

			// Only SSH URL provided.
			Entry(
				`[x] ssh [ ] http -- [ ] ssh agent [ ] ssh key [ ] prefer http`,
				true, false,
				false, false, false,
				"SSH is the only available protocol but there is no SSH agent and no private key was provided",
			),
			Entry(
				`[x] ssh [ ] http -- [ ] ssh agent [ ] ssh key [x] prefer http`,
				true, false,
				false, false, true,
				"SSH is the only available protocol but there is no SSH agent and no private key was provided",
			),
			Entry(
				`[x] ssh [ ] http -- [ ] ssh agent [x] ssh key [ ] prefer http`,
				true, false,
				false, true, false,
				"ssh",
			),
			Entry(
				`[x] ssh [ ] http -- [ ] ssh agent [x] ssh key [x] prefer http`,
				true, false,
				false, true, true,
				"ssh",
			),
			Entry(
				`[x] ssh [ ] http -- [x] ssh agent [ ] ssh key [ ] prefer http`,
				true, false,
				true, false, false,
				"ssh",
			),
			Entry(
				`[x] ssh [ ] http -- [x] ssh agent [ ] ssh key [x] prefer http`,
				true, false,
				true, false, true,
				"ssh",
			),
			Entry(
				`[x] ssh [ ] http -- [x] ssh agent [x] ssh key [ ] prefer http`,
				true, false,
				true, true, false,
				"ssh",
			),
			Entry(
				`[x] ssh [ ] http -- [x] ssh agent [x] ssh key [x] prefer http`,
				true, false,
				true, true, true,
				"ssh",
			),

			// Both URLs provided.
			Entry(
				`[x] ssh [x] http -- [ ] ssh agent [ ] ssh key [ ] prefer http`,
				true, true,
				false, false, false,
				"http",
			),
			Entry(
				`[x] ssh [x] http -- [ ] ssh agent [ ] ssh key [x] prefer http`,
				true, true,
				false, false, true,
				"http",
			),
			Entry(
				`[x] ssh [x] http -- [ ] ssh agent [x] ssh key [ ] prefer http`,
				true, true,
				false, true, false,
				"ssh",
			),
			Entry(
				`[x] ssh [x] http -- [ ] ssh agent [x] ssh key [x] prefer http`,
				true, true,
				false, true, true,
				"http",
			),
			Entry(
				`[x] ssh [x] http -- [x] ssh agent [ ] ssh key [ ] prefer http`,
				true, true,
				true, false, false,
				"ssh",
			),
			Entry(
				`[x] ssh [x] http -- [x] ssh agent [ ] ssh key [x] prefer http`,
				true, true,
				true, false, true,
				"http",
			),
			Entry(
				`[x] ssh [x] http -- [x] ssh agent [x] ssh key [ ] prefer http`,
				true, true,
				true, true, false,
				"ssh",
			),
			Entry(
				`[x] ssh [x] http -- [x] ssh agent [x] ssh key [x] prefer http`,
				true, true,
				true, true, true,
				"http",
			),
		)
	})
})

// expectCloneWithURL expects a local Git clone to exist in the given directory,
// with the origin remote using the given URL.
//
// The logger is inspected to verify it contains the output from Git itself.
func expectCloneWithURL(dir, url string, logger *logging.BufferedLogger) {
	repo, err := git.PlainOpen(dir)
	Expect(err).ShouldNot(HaveOccurred())

	rem, err := repo.Remote("origin")
	Expect(err).ShouldNot(HaveOccurred())
	Expect(rem.Config().URLs).To(ConsistOf(url))
	Expect(logger.Messages()).To(ContainElement(
		logging.BufferedLogMessage{
			Message: "git: Total 3 (delta 0), reused 3 (delta 0), pack-reused 0",
		},
	))
}
