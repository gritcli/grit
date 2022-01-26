package flags

import (
	"github.com/gritcli/grit/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// SetupSocket sets up the --socket flag on the root command.
func SetupSocket(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"socket",
		config.DefaultDaemonSocket,
		"set the Unix socket address of the Grit daemon",
	)
	cmd.MarkPersistentFlagFilename("socket")
}

// Socket returns the path to the Unix socket address of the Grit daemon.
func Socket(cmd *cobra.Command) (string, error) {
	dir, err := cmd.Flags().GetString("socket")
	if err != nil {
		panic(err)
	}

	return homedir.Expand(dir)
}
