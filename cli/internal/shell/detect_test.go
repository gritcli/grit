package shell_test

import (
	"os"
	"strings"

	. "github.com/gritcli/grit/cli/internal/shell"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Detect()", func() {
	DescribeTable(
		"it detects the shell type from environment variables",
		func(expect Type, env ...string) {
			snapshot := os.Environ()
			setEnv(env)
			defer setEnv(snapshot)

			Expect(Detect()).To(Equal(expect))
		},
		Entry(
			"unknown",
			UnknownType,
		),
		Entry(
			"zsh",
			ZshType,
			"SHELL=/usr/bin/zsh",
		),
		Entry(
			"bash",
			BashType,
			"SHELL=/bin/bash",
		),
		Entry(
			"fish",
			FishType,
			"SHELL=/usr/local/bin/fish",
		),
		Entry(
			"powershell",
			PowerShellType,
			"PSVersionTable=...",
		),
	)
})

// setEnv sets the entire environment to match the key/value pairs in env.
func setEnv(env []string) {
	os.Clearenv()

	for _, kv := range env {
		pair := strings.SplitN(kv, "=", 2)
		os.Setenv(pair[0], pair[1])
	}
}
