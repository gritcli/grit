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
	"github.com/gritcli/grit/common/config"
	"golang.org/x/oauth2"
)

// Driver is an implementation of source.Driver that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type Driver struct {
	Config config.GitHub
	Logger logging.Logger

	client *github.Client
	cache  cache
}

// Init initializes the driver.
func (d *Driver) Init(ctx context.Context) error {
	httpClient := http.DefaultClient
	if d.Config.Token != "" {
		httpClient = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: d.Config.Token},
			),
		)
	}

	if isGitHubDotCom(d.Config) {
		d.client = github.NewClient(httpClient)
	} else {
		var err error
		d.client, err = github.NewEnterpriseClient(d.Config.Domain, "", httpClient)
		if err != nil {
			return err
		}
	}

	user, res, err := d.client.Users.Get(ctx, "")
	if err != nil {
		if res.StatusCode != http.StatusUnauthorized {
			return err
		}

		logging.Log(d.Logger, "not authenticated")
		return nil
	}

	logging.Log(d.Logger, "authenticated as %s", user.GetLogin())
	d.cache.SetCurrentUser(user)

	if err := d.populateRepoCache(ctx); err != nil {
		return err
	}

	return nil
}

// Run performs any ongoing behavior required by the driver.
func (d *Driver) Run(ctx context.Context) error {
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

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func (d *Driver) populateRepoCache(ctx context.Context) error {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var repos []*github.Repository

	for opts.Page != 0 {
		repoPage, res, err := d.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, r := range repoPage {
			logging.Debug(d.Logger, "cached repository: %s", r.GetFullName())
			repos = append(repos, r)
		}

		opts.Page = res.NextPage
	}

	logging.Log(
		d.Logger,
		"cached %d repositories",
		len(repos),
	)

	d.cache.SetRepos(repos)

	return nil
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(cfg config.GitHub) bool {
	return strings.EqualFold(cfg.Domain, "github.com")
}
