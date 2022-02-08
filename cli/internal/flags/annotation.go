package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Annotation contains grit-specific annotations attached to flags.
type Annotation struct {
	ArgName string
}

// Annotate adds an annotation to a flag.
func Annotate(cmd *cobra.Command, name string, ann Annotation) {
	f := cmd.Flags().Lookup(name)
	if f == nil {
		panic("unknown flag")
	}

	if f.Annotations == nil {
		f.Annotations = map[string][]string{}
	}

	f.Annotations["grit"] = []string{
		ann.ArgName,
	}
}

// GetAnnotation returns the annotation on a given flag.
func GetAnnotation(f *pflag.Flag) Annotation {
	if ann, ok := f.Annotations["grit"]; ok {
		return Annotation{
			ann[0],
		}
	}

	return Annotation{}
}
