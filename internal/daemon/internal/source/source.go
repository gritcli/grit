package source

import "github.com/gritcli/grit/plugin/sourcedriver"

// Source is a repository source.
type Source struct {
	// Name is the unique name for the repository source.
	Name string

	// Description is a human-readable description of the source.
	Description string

	// CloneDir is the directory containing repositories cloned from this
	// source.
	CloneDir string

	// Driver is the driver used to perform repository operations.
	Driver sourcedriver.Driver
}
