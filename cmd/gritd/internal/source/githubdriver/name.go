package githubdriver

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// ownerNamePattern is a regex that matches valid GitHub "owner" names (such
	// as usernames and organization names).
	ownerNamePattern = regexp.MustCompile(`(?i)^[a-z0-9]+(?:-[a-z0-9]+)*$`)

	// repoNamePattern is a regex that matches valid GitHub repository names.
	repoNamePattern = regexp.MustCompile(`(?i)^[a-z0-9_\-\.]+$`)
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
	repoName = name
	if i := strings.IndexRune(name, '/'); i > 0 {
		ownerName = name[:i]
		repoName = name[i+1:]

		if !ownerNamePattern.MatchString(ownerName) {
			return "", "", fmt.Errorf("repository name (%s) contains an invalid owner component", name)
		}
	}

	if !repoNamePattern.MatchString(repoName) {
		return "", "", fmt.Errorf("repository name (%s) contains an invalid repository component", name)
	}

	return ownerName, repoName, nil
}
