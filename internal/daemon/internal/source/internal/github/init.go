package github

import (
	"context"
	"net/http"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Init initializes the driver.
func (d *Driver) Init(
	ctx context.Context,
	logger logging.Logger,
) error {
	httpClient := http.DefaultClient
	if d.Config.Token != "" {
		httpClient = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: d.Config.Token},
			),
		)
	}

	if isEnterpriseServer(d.Config.Domain) {
		var err error
		d.client, err = github.NewEnterpriseClient(d.Config.Domain, "", httpClient)
		if err != nil {
			return err
		}
	} else {
		d.client = github.NewClient(httpClient)
	}

	if d.Config.Token == "" {
		logging.Log(logger, "not authenticated (no token specified)")
		return nil
	}

	user, res, err := d.client.Users.Get(ctx, "")
	if err != nil {
		if res.StatusCode != http.StatusUnauthorized {
			return err
		}

		// TODO: rebuild client without token provider
		logging.Log(logger, "not authenticated (token is invalid)")
		return nil
	}

	logging.Log(logger, "authenticated as %s", user.GetLogin())
	d.cache.SetCurrentUser(user)

	if err := d.populateRepoCache(ctx, logger); err != nil {
		return err
	}

	return nil
}

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func (d *Driver) populateRepoCache(
	ctx context.Context,
	logger logging.Logger,
) error {
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
			logging.Debug(logger, "discovered %s", r.GetFullName())
			repos = append(repos, r)
		}

		opts.Page = res.NextPage
	}

	logging.Log(
		logger,
		"added %d repositories to the repository list for @%s",
		len(repos),
		d.cache.CurrentUser().GetLogin(),
	)

	d.cache.SetRepos(repos)

	return nil
}
