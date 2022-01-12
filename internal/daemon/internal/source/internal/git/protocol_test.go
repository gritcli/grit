package git

import (
	"os"

	"github.com/gritcli/grit/internal/daemon/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func useHTTP()", func() {
	DescribeTable(
		"it chooses the best available protocol",
		func(
			hasSSH, hasHTTP,
			hasSSHAgent, hasSSHPrivateKey, preferHTTP bool,
			expect string, // "ssh", "http" or an error message
		) {
			cfg := config.Git{
				PreferHTTP: preferHTTP,
			}

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
				cfg.SSHKeyFile = "<not empty>"
			}

			useHTTP, err := useHTTP(hasSSH, hasHTTP, cfg)

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
