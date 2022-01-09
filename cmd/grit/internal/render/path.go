package render

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// AbsPath formats file or directory path, substituting ~ for the current user's
// home directory, if possible.
func AbsPath(p string) string {
	p = filepath.Clean(p)

	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}

	if h, ok := compactHomeDir(abs); ok {
		return h
	}

	return p
}

// RelPath formats file or directory path in the shortest way possible, using
// paths relative to the current working directory, if possible.
func RelPath(p string) string {
	shortest := filepath.Clean(p)

	abs, err := filepath.Abs(p)
	if err == nil && len(abs) < len(shortest) {
		shortest = abs
	}

	base, err := os.Getwd()
	if err == nil {
		rel, err := filepath.Rel(base, abs)
		if err == nil && len(rel) < len(shortest) {
			shortest = rel
		}
	}

	if h, ok := compactHomeDir(abs); ok && len(h) < len(shortest) {
		shortest = h
	}

	return shortest
}

// compactHomeDir returns p compacted with ~ syntax if p is within the current
// user's home directory.
func compactHomeDir(p string) (string, bool) {
	base, err := homedir.Dir()
	if err != nil {
		return p, false
	}

	rel, err := filepath.Rel(base, p)
	if err != nil {
		return p, false
	}

	return filepath.Join("~", rel), true
}
