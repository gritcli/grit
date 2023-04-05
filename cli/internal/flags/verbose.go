package flags

import "github.com/spf13/cobra"

// SetupVerbose sets up the --verbose flag on the root command.
func SetupVerbose(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
}

// IsVerbose returns true if the --verbose flag was set.
func IsVerbose(cmd *cobra.Command) bool {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		panic(err)
	}

	return verbose
}
