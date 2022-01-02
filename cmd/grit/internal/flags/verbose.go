package flags

import "github.com/spf13/cobra"

// SetupVerbose sets up the --verbose flag on the root command.
func SetupVerbose(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
}

// IsVerbose checks if --verbose was provided.
func IsVerbose(cmd *cobra.Command) bool {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		panic(err)
	}

	return verbose
}
