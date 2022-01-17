package sourcedriver

// RemoteRepo is a reference to a remote repository provided by a source.
type RemoteRepo struct {
	// ID uniquely identifies the repository within the source.
	ID string

	// Name is the human-readable name of the repository.
	//
	// The formatting and uniqueness guarantees of the repository name are
	// driver-specific.
	Name string

	// Description is a human-readable description of the repository.
	Description string

	// WebURL is the URL of the repository's web page, if available.
	//
	// This is a page viewable in a browser by a human, not the URL used to
	// clone the repository.
	WebURL string
}
