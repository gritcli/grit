package githubdriver

import (
	"errors"
	"fmt"
	"strings"
)

// parseRepoName parses a repository name into its owner and unqualified name
// components.
//
// If the name is fully-qualified (contains a slash), then ownerName is the part
// before the slash and repoName is the part after the slash.
//
// if the name is NOT fully-qualified (does not contain a slash) then ownerName
// is empty and repoName is equal to name.
func parseRepoName(name string) (ownerName, repoName string, err error) {
	if name == "" {
		return "", "", errors.New("repository name must not be empty")
	}

	i := strings.IndexRune(name, '/')
	if i == -1 {
		return "", name, nil
	}

	ownerName = name[:i]
	repoName = name[i+1:]

	if ownerName == "" {
		return "", "", fmt.Errorf("repository name (%s) contains an empty owner component", name)
	}

	if repoName == "" {
		return "", "", fmt.Errorf("repository name (%s) contains an empty name component", name)
	}

	return ownerName, repoName, nil
}
