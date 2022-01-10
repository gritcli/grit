package scm

import "context"

// Cloner is an interface for cloning a repository.
//
// The term "clone" is intended in the manner used by Git and similar
// distributed source control systems. The implementation should take
// whatever action best matches this intent. For example, for Subversion
// repositories an equivalent operation would be "checkout".
type Cloner interface {
	// Clone clones the repository into the given target directory.
	Clone(ctx context.Context, dir string) error
}
