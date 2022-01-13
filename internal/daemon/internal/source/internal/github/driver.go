package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dogmatiq/dodeca/logging"
	humanize "github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/daemon/internal/config"
)

// Driver is an implementation of driver.Driver that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type Driver struct {
	Config config.GitHub

	client *github.Client
	cache  cache
}

// Run performs any ongoing behavior required by the driver.
func (d *Driver) Run(
	ctx context.Context,
	logger logging.Logger,
) error {
	return nil
}

// Status returns a brief description of the status of the driver.
func (d *Driver) Status(ctx context.Context) (string, error) {
	invalidToken := false
	limits, _, err := d.client.RateLimits(ctx)
	if err != nil {
		var e *github.ErrorResponse

		if errors.As(err, &e) {
			if e.Response.StatusCode != http.StatusUnauthorized {
				return "", err
			}

			// This endpoint does not require authentication, so we can only get
			// an unauthorized error if we explicitly provided invalid
			// credentials.
			invalidToken = true
		} else {
			return "", err
		}
	}

	var info []string

	if invalidToken {
		info = append(info, "unauthenticated (invalid token)")
	} else {
		if u := d.cache.CurrentUser(); u != nil {
			info = append(info, "@"+u.GetLogin())
		} else {
			info = append(info, "unauthenticated")
		}

		info = append(
			info,
			fmt.Sprintf(
				"%d API requests remaining (resets %s)",
				limits.GetCore().Remaining,
				humanize.Time(
					limits.GetCore().Reset.Time,
				),
			),
		)
	}

	return strings.Join(info, ", "), nil
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(cfg config.GitHub) bool {
	return strings.EqualFold(cfg.Domain, "github.com")
}
