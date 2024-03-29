package flags

import (
	"github.com/spf13/cobra"
)

// SetupShellExecutorOutput sets up --shell-executor-output flag on the root command.
func SetupShellExecutorOutput(cmd *cobra.Command) {
	// Add --shell-executor-output as a "persistent" flag so it's available on
	// all commands. Mark it as hidden as it should only be passed by the code
	// generated by the "setup-shell" command, and never by the user directly.
	f := cmd.PersistentFlags()
	f.String("shell-executor-output", "", "output file for shell commands to execute")

	if err := f.MarkHidden("shell-executor-output"); err != nil {
		panic(err)
	}
}

// ShellExecutorOutputFile returns the path to the file that the shell.Executor
// should write its command to, if any.
func ShellExecutorOutputFile(cmd *cobra.Command) (string, bool) {
	filename, err := cmd.Flags().GetString("shell-executor-output")
	if err != nil {
		panic(err)
	}

	return filename, filename != ""
}
