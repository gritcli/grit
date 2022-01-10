package source

import "strings"

// List is a collection of sources.
type List []Source

// ByName returns the source with the given name.
//
// Source names are case insensitive.
func (l List) ByName(n string) (Source, bool) {
	for _, src := range l {
		if strings.EqualFold(src.Name, n) {
			return src, true
		}
	}

	return Source{}, false
}
