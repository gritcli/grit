package flags

import (
	"io"

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

	return nonInteractive
}

// supportsInteractivity returns true if cmd can support interactive prompts.
func supportsInteractivity(cmd *cobra.Command) bool {
	// The promptui library requires close capabilities on the IO reader/writer
	// for some reason, so we can't perform any interactivity unless this
	// operation is available.
	if _, ok := cmd.InOrStdin().(io.ReadCloser); !ok {
		return false
	}

	if _, ok := cmd.OutOrStdout().(io.WriteCloser); !ok {
		return false
	}

	// Lastly, we want to check if the file we're writing to is actually a
	// terminal.
	type fileDescriptor interface {
		Fd() uintptr
	}

	if fd, ok := cmd.OutOrStdout().(fileDescriptor); ok {
		return isatty.IsTerminal(fd.Fd())
	}

	return false
}
