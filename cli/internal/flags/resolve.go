package flags

import (
	"errors"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/completion"
	"github.com/spf13/cobra"
)

// SetupFromSource sets up the --from-source (and --no-resolve) flags used by
// commands that resolve query strings to repositories.
func SetupFromSource(cmd *cobra.Command, con *imbue.Container) {
	cmd.Flags().StringP(
		"from-source", "f",
		"",
		"limit resolution of <repo> to a single `source`",
	)

	cmd.Flags().Bool(
		"no-resolve",
		false,
		"force <repo> to be interpreted as a unique ID, requires --from-source",
	)

	cmd.RegisterFlagCompletionFunc(
		"from-source",
		completion.SourceName(con),
	)
}

// FromSource returns the source name passed via --from-source. It returns an
// empty string if --from-source is omitted.
func FromSource(cmd *cobra.Command) (source string, noResolve bool, _ error) {
	source, err := cmd.Flags().GetString("from-source")
	if err != nil {
		panic(err)
	}

	noResolve, err = cmd.Flags().GetBool("no-resolve")
	if err != nil {
		panic(err)
	}

	if noResolve && source == "" {
		return "", false, errors.New("--no-resolve requires --from-source")
	}

	return source, noResolve, nil
}
