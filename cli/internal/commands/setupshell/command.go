package setupshell

import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

// Command returns the "setup-shell" command.
func Command(con *imbue.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "setup-shell",
		DisableFlagsInUseLine: true,
		Short:                 "Generate auto-completion and shell integration code",
	}

	for _, sh := range supportedShells() {
		sh := sh // capture loop variable

		cmd.AddCommand(
			&cobra.Command{
				Use:                   sh,
				DisableFlagsInUseLine: true,
				Short:                 fmt.Sprintf("Generate integration code for %s", sh),
				RunE: func(cmd *cobra.Command, args []string) error {
					f, err := installers.Open("install." + sh)
					if err != nil {
						return err
					}
					defer f.Close()

					code, err := io.ReadAll(f)
					if err != nil {
						return err
					}

					bin, err := os.Executable()
					if err != nil {
						return err
					}

					cmd.Print(
						strings.ReplaceAll(
							string(code),
							"/path/to/grit",
							shell.Escape(bin),
						),
					)

					return err
				},
			},
		)
	}

	return cmd
}

//go:embed install.*
var installers embed.FS

// supportedShells returns the list of shells that have installers.
func supportedShells() []string {
	files, err := installers.ReadDir(".")
	if err != nil {
		panic(err)
	}

	var shells []string
	for _, file := range files {
		shells = append(
			shells,
			strings.TrimPrefix(
				filepath.Ext(
					file.Name(),
				),
				".",
			),
		)
	}
	sort.Strings(shells)

	return shells
}
