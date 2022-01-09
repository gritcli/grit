package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dogmatiq/dodeca/logging"
	humanize "github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/config"
	"golang.org/x/oauth2"
)

// driver is an implementation of source.Driver that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type driver struct {
	cfg    config.GitHub
	client *github.Client
	logger logging.Logger

	cache cache
}

// NewDriver returns a new driver with the given configuration.
func NewDriver(
	cfg config.GitHub,
	logger logging.Logger,
) (source.Driver, error) {
	d := &driver{
		cfg:    cfg,
		logger: logger,
	}

	httpClient := http.DefaultClient
	if cfg.Token != "" {
		httpClient = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: cfg.Token},
			),
		)
	}

	if isGitHubDotCom(cfg) {
		d.client = github.NewClient(httpClient)
	} else {
		var err error
		d.client, err = github.NewEnterpriseClient(cfg.Domain, "", httpClient)
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

// Init initializes the driver.
func (d *driver) Init(ctx context.Context) error {
	user, res, err := d.client.Users.Get(ctx, "")
	if err != nil {
		if res.StatusCode != http.StatusUnauthorized {
			return err
		}

		logging.Log(d.logger, "not authenticated")
		return nil
	}

	logging.Log(d.logger, "authenticated as %s", user.GetLogin())
	d.cache.SetCurrentUser(user)

	if err := d.populateRepoCache(ctx); err != nil {
		return err
	}

	return nil
}

// Run performs any ongoing behavior required by the driver.
func (d *driver) Run(ctx context.Context) error {
	return nil
}

// Status returns a brief description of the status of the driver.
func (d *driver) Status(ctx context.Context) (string, error) {
	limits, res, err := d.client.RateLimits(ctx)
	if err != nil {
		if res == nil || res.StatusCode != http.StatusUnauthorized {
			return "", err
		}
	}

	var info []string

	if res.StatusCode == http.StatusUnauthorized {
		info = append(info, "unauthenticated (bad credentials)")
	} else if u := d.cache.CurrentUser(); u != nil {
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

	return strings.Join(info, ", "), nil
}

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func (d *driver) populateRepoCache(ctx context.Context) error {
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
			logging.Debug(d.logger, "cached repository: %s", r.GetFullName())
			repos = append(repos, r)
		}

		opts.Page = res.NextPage
	}

	logging.Log(
		d.logger,
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
