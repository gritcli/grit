package commands

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	"github.com/charmbracelet/glamour"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/commands/clone"
	"github.com/gritcli/grit/cli/internal/commands/source"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

// Root is the root "grit" command.
var Root = &cobra.Command{
	Use:   "grit",
	Short: "Manage your local VCS clones",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Add the currently-executing Cobra CLI command the the DI
		// container.
		//
		// This hook is called after the CLI arguments are resolved to a
		// specific command.
		//
		// This allows other DI provider definitions to make use of the
		// flags passed to the command.
		cobradi.Provide(cmd, func() *cobra.Command {
			return cmd
		})
	},
}

var (
	//go:embed help.gtpl
	helpTemplate string

	// go:embed usage.gtpl
	usageTemplate string
)

func init() {
	flags.SetupVerbose(Root)
	flags.SetupNoInteractive(Root)
	flags.SetupSocket(Root)
	flags.SetupShellExecutorOutput(Root)

	Root.SetHelpFunc(help)

	// Root.SetHelpTemplate(helpTemplate)
	// Root.SetUsageTemplate(usageTemplate)

	Root.AddCommand(
		clone.Command,
		source.Command,
		&cobra.Command{
			Use:  "additional",
			Long: "what",
		},
		generateDocs,
	)
}

var generateDocs = &cobra.Command{
	Use:    "generate-docs",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return doc.GenManTree(cmd.Root(), &doc.GenManHeader{
			// Title:   "",
			// Section: "",
			// Date:    &time.Time{},
			// Source:  "",
			// Manual:  "",
		}, "/tmp/manpage")
	},
}

type flag struct {
	Short   string
	Long    string
	Arg     string
	Usage   string
	Default string
}

func marshalFlag(in *pflag.Flag) flag {
	ann := flags.GetAnnotation(in)

	out := flag{
		Short: in.Shorthand,
		Long:  in.Name,
		Arg:   ann.ArgName,
		Usage: in.Usage,
	}

	fmt.Printf("%#v\n", in)

	if in.ShorthandDeprecated != "" {
		out.Short = ""
	}

	if in.DefValue != "" {
		if in.Value.Type() != "bool" || in.DefValue != "false" {
			out.Default = in.DefValue
		}
	}

	return out
}

func help(cmd *cobra.Command, args []string) {
	ctx := struct {
		*cobra.Command
		LocalFlags  []flag
		GlobalFlags []flag
	}{
		Command: cmd,
	}

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			ctx.LocalFlags = append(ctx.LocalFlags, marshalFlag(f))
		}
	})

	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			ctx.GlobalFlags = append(ctx.GlobalFlags, marshalFlag(f))
		}
	})

	t := template.New("help")
	template.Must(t.Parse(helpTemplate))

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, ctx); err != nil {
		panic(err)
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("auto"),
		// glamour.WithWordWrap(-1),
	)
	if err != nil {
		panic(err)
	}

	out, err := r.RenderBytes(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if _, err := cmd.OutOrStdout().Write(out); err != nil {
		panic(err)
	}
}

// 	line := ""
// 	if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
// 		line = fmt.Sprintf("  -%s, --%s", flag.Shorthand, flag.Name)
// 	} else {
// 		line = fmt.Sprintf("      --%s", flag.Name)
// 	}

// 	varname, usage := UnquoteUsage(flag)
// 	if varname != "" {
// 		line += " " + varname
// 	}
// 	if flag.NoOptDefVal != "" {
// 		switch flag.Value.Type() {
// 		case "string":
// 			line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
// 		case "bool":
// 			if flag.NoOptDefVal != "true" {
// 				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
// 			}
// 		case "count":
// 			if flag.NoOptDefVal != "+1" {
// 				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
// 			}
// 		default:
// 			line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
// 		}
// 	}

// 	// This special character will be replaced with spacing once the
// 	// correct alignment is calculated
// 	line += "\x00"
// 	if len(line) > maxlen {
// 		maxlen = len(line)
// 	}

// 	line += usage
// 	if !flag.defaultIsZeroValue() {
// 		if flag.Value.Type() == "string" {
// 			line += fmt.Sprintf(" (default %q)", flag.DefValue)
// 		} else {
// 			line += fmt.Sprintf(" (default %s)", flag.DefValue)
// 		}
// 	}
// 	if len(flag.Deprecated) != 0 {
// 		line += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
// 	}

// 	lines = append(lines, line)
