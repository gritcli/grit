package githubsource

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v38/github"
)

// Source is an implementation of source.Source for GitHub servers.
type Source struct {
	Client *github.Client
}

// Description returns a short, human-readable description of the source.
//
// The description should be adequate to distinguish this source from any
// other sources that may exist.
func (s *Source) Description() string {
	if s.isDotCom() {
		return "github"
	}

	return fmt.Sprintf(
		"github enterprise %s",
		s.Client.BaseURL.Host,
	)
}

// Status queries the status of the source.
//
// It returns an error if the source is misconfigured or unreachable.
//
// The status string should include any source-specific information
func (s *Source) Status(ctx context.Context) (string, error) {
	u, scopes, ok, err := s.queryUser(ctx)
	if err != nil {
		return "", err
	}

	if !ok {
		return "unauthenticated", nil
	}

	status := fmt.Sprintf("@%s", u.GetLogin())

	for _, s := range diffScopes(RequiredScopes, scopes) {
		status += fmt.Sprintf(", missing '%s' scope", s)
	}

	for _, s := range diffScopes(scopes, RequiredScopes) {
		status += fmt.Sprintf(", unnecessary '%s' scope", s)
	}

	return status, nil
}

// isDotCom returns true if this source is used for github.com, as opposed to a
// GitHub Enterprise server.
func (s *Source) isDotCom() bool {
	return s.Client.BaseURL.Host == "api.github.com"
}

// queryUser queries the currently authenticated user, or returns false if not
// authenticated.
func (s *Source) queryUser(ctx context.Context) (*github.User, []string, bool, error) {
	user, resp, err := s.Client.Users.Get(ctx, "")
	if err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, nil, false, nil
		}

		return nil, nil, false, err
	}

	scopes := strings.Fields(
		resp.Header.Get("X-OAuth-Scopes"),
	)

	return user, scopes, true, nil
}
