package flags

import (
	"github.com/gritcli/grit/internal/common/api"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// SetupSocket sets up the --socket flag on the root command.
func SetupSocket(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"socket",
		api.DefaultSocket,
		"set the Unix socket address of the Grit daemon",
	)
}

// Socket returns the path to the Unix socket address of the Grit daemon.
func Socket(cmd *cobra.Command) (string, error) {
	dir, err := cmd.Flags().GetString("socket")
	if err != nil {
		panic(err)
	}

	return homedir.Expand(dir)
}
