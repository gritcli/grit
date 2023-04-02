package gitvcs // note: no _test suffix to allow testing unexported useHTTP() method.

import (
	"context"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/gritcli/grit/logs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Cloner", func() {
	var (
		ctx     context.Context
		buffer  logs.Buffer
		cloner  *Cloner
		tempDir string
	)

	BeforeEach(func() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		DeferCleanup(cancel)

		buffer = logs.Buffer{}
		cloner = &Cloner{}

		var err error
		tempDir, err = os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(func() {
			os.RemoveAll(tempDir)
		})
	})

	Describe("func Clone()", func() {
		// Note, in order to minimise the number of clones, these tests do not
		// cover the full range of configuration options available. These are
		// covered in the tests for the useHTTP() method.

		It("clones via SSH using the SSH agent", func() {
			if os.Getenv("SSH_AUTH_SOCK") == "" {
				Skip("SSH agent is not available")
			}

			// Use the public test repo; it makes it easier to run the tests
			// locally using the developer's SSH agent.
			cloner.SSHEndpoint = "git@github.com:grit-integration-tests/test-public.git"

			err := cloner.Clone(ctx, tempDir, buffer.Log())
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.SSHEndpoint, buffer)
		})

		It("clones via SSH using an explicit private key", func() {
			// Use the private test repo to ensure the private key was
			// definitely used.
			cloner.SSHEndpoint = "git@github.com:grit-integration-tests/test-private.git"
			cloner.SSHKeyFile = "./testdata/deploy-key-no-passphrase"

			err := cloner.Clone(ctx, tempDir, buffer.Log())
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.SSHEndpoint, buffer)
		})

		It("clones via SSH using an explicit private key with a passphrase", func() {
			// Use the private test repo to ensure the private key was
			// definitely used.
			cloner.SSHEndpoint = "git@github.com:grit-integration-tests/test-private.git"
			cloner.SSHKeyFile = "./testdata/deploy-key-with-passphrase"
			cloner.SSHKeyPassphrase = "passphrase"

			err := cloner.Clone(ctx, tempDir, buffer.Log())
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.SSHEndpoint, buffer)
		})

		It("clones via HTTP without authentication", func() {
			// Use the public test repo so that it's accessible without
			// authentication.
			cloner.HTTPEndpoint = "https://github.com/grit-integration-tests/test-public.git"
			cloner.PreferHTTP = true

			err := cloner.Clone(ctx, tempDir, buffer.Log())
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.HTTPEndpoint, buffer)
		})

		It("clones via HTTP with authentication", func() {
			token := os.Getenv("GRIT_INTEGRATION_TEST_GITHUB_TOKEN")
			if token == "" {
				Skip("GRIT_INTEGRATION_TEST_GITHUB_TOKEN is not set")
			}

			// Use the private test repo to ensure the basic authentication
			// details were actually used.
			cloner.HTTPEndpoint = "https://github.com/grit-integration-tests/test-private.git"
			cloner.HTTPUsername = "<ignored-by-github>"
			cloner.HTTPPassword = token
			cloner.PreferHTTP = true

			err := cloner.Clone(ctx, tempDir, buffer.Log())
			Expect(err).ShouldNot(HaveOccurred())
			expectCloneWithURL(tempDir, cloner.HTTPEndpoint, buffer)
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
				if hasSSH {
					cloner.SSHEndpoint = "<ssh endpoint>"
				}

				if hasHTTP {
					cloner.HTTPEndpoint = "<http endpoint>"
				}

				cloner.PreferHTTP = preferHTTP

				orig := os.Getenv("SSH_AUTH_SOCK")
				defer os.Setenv("SSH_AUTH_SOCK", orig)

				if hasSSHAgent {
					os.Setenv("SSH_AUTH_SOCK", "<ssh auth socket>")
				} else {
					os.Setenv("SSH_AUTH_SOCK", "")
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
// The log buffer is inspected to verify it contains the output from Git itself.
func expectCloneWithURL(dir, url string, buffer logs.Buffer) {
	repo, err := git.PlainOpen(dir)
	Expect(err).ShouldNot(HaveOccurred())

	rem, err := repo.Remote("origin")
	Expect(err).ShouldNot(HaveOccurred())
	Expect(rem.Config().URLs).To(ConsistOf(url))

	Expect(buffer).NotTo(BeEmpty())
	Expect(buffer[len(buffer)-1].Text).To(
		MatchRegexp(`git: Total \d+ \(delta \d+\), reused \d+ \(delta \d+\), pack-reused \d+`),
	)
}
