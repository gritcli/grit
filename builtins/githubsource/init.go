package githubsource

import (
	"context"
	"net/http"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Init initializes the driver.
func (d *impl) Init(
	ctx context.Context,
	logger logging.Logger,
) error {
	httpClient := http.DefaultClient
	if d.config.Token != "" {
		httpClient = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: d.config.Token},
			),
		)
	}

	if isEnterpriseServer(d.config.Domain) {
		var err error
		d.client, err = github.NewEnterpriseClient(d.config.Domain, "", httpClient)
		if err != nil {
			return err
		}
	} else {
		d.client = github.NewClient(httpClient)
	}

	if d.config.Token == "" {
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
	d.user = user

	if err := d.populateRepoCache(ctx, logger); err != nil {
		return err
	}

	return nil
}

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func (d *impl) populateRepoCache(
	ctx context.Context,
	logger logging.Logger,
) error {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	d.reposByID = map[int64]*github.Repository{}
	d.reposByOwner = map[string]map[string]*github.Repository{}

	for opts.Page != 0 {
		repoPage, res, err := d.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, r := range repoPage {
			logging.Debug(logger, "discovered %s", r.GetFullName())

			owner := r.GetOwner().GetLogin()
			reposByName := d.reposByOwner[owner]
			if reposByName == nil {
				reposByName = map[string]*github.Repository{}
				d.reposByOwner[owner] = reposByName
			}

			reposByName[r.GetName()] = r
			d.reposByID[r.GetID()] = r
		}

		opts.Page = res.NextPage
	}

	logging.Log(
		logger,
		"added %d repositories to the repository list for @%s",
		len(d.reposByID),
		d.user.GetLogin(),
	)

	return nil
}
