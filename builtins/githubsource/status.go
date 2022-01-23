package githubsource

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
)

// Status returns a brief description of the status of the driver.
func (d *impl) Status(ctx context.Context) (string, error) {
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
		if u := d.user; u != nil {
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
