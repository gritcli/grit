package flags

import (
	"errors"

	"github.com/spf13/cobra"
)

// SetupFromSource sets up the --from-source (and --no-resolve) flags used by
// commands that resolve query strings to repositories.
func SetupFromSource(cmd *cobra.Command) {
	cmd.Flags().StringP(
		"from-source", "f",
		"asdfas",
		"Limit resolution of `<repo>` to a single source. By default all sources are searched.",
	)
	Annotate(cmd, "from-source", Annotation{
		ArgName: "source",
	})

	// 	`--from-source|-f <source>`
	//   Limit resolution of <repo> to the given source. By default all
	//   sources are searched.

	// `--no-resolve`
	//   Force `<repo>` to be interpreted as a unique identifier.
	//   Requires `--source`.

	cmd.Flags().Bool(
		"no-resolve",
		false,
		"Force `<repo>` to be interpreted as a (source-specific) unique identifier, `requires --from-source`",
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
