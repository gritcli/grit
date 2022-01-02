package github

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
)

// convertRepo converts a github.Repository to a source.Repo.
func convertRepo(r *github.Repository) source.Repo {
	return source.Repo{
		ID:          strconv.FormatInt(r.GetID(), 10),
		Name:        r.GetFullName(),
		Description: r.GetDescription(),
		WebURL:      r.GetHTMLURL(),
	}
}

// parseRepoID parses a repo ID from its string form (as used by the source
// package) to the numeric form used by the GitHub API.
func parseRepoID(id string) (int64, error) {
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid repo ID: %w", err)
	}

	if n <= 0 {
		return 0, errors.New("invalid repo ID: expected positive integer")
	}

	return n, nil
}

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
