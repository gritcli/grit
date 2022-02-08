package flags

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// SetupNoInteractive sets up the --no-interaction flag on the root command.
func SetupNoInteractive(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(
		"no-interaction", "n",
		!supportsInteractivity(cmd), // default to true if cmd doesn't appear to support interactivity
		"do not ask any interactive questions",
	)
}

// IsInteractive returns true if cmd allows interactive prompts.
func IsInteractive(cmd *cobra.Command) bool {
	nonInteractive, err := cmd.Flags().GetBool("no-interaction")
	if err != nil {
		panic(err)
	}

	return !nonInteractive
}

// supportsInteractivity returns true if cmd can support interactive prompts.
func supportsInteractivity(cmd *cobra.Command) bool {
	return isTerminal(cmd.InOrStdin()) &&
		isTerminal(cmd.OutOrStdout())
}

// isTerminal returns true if x is an *os.File that is a TTY.
func isTerminal(x interface{}) bool {
	if f, ok := x.(*os.File); ok {
		return isatty.IsTerminal(f.Fd())
	}

	return false
}
