package githubsource

import (
	"context"
	"net/http"

	"github.com/google/go-github/v50/github"
	"github.com/gritcli/grit/logs"
	"golang.org/x/oauth2"
)

// Init initializes the source.
func (s *source) Init(
	ctx context.Context,
	log logs.Log,
) error {
	httpClient := http.DefaultClient
	if s.config.Token != "" {
		httpClient = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: s.config.Token},
			),
		)
	}

	if isEnterpriseServer(s.config.Domain) {
		var err error
		s.client, err = github.NewEnterpriseClient(s.config.Domain, "", httpClient)
		if err != nil {
			return err
		}
	} else {
		s.client = github.NewClient(httpClient)
	}

	if s.config.Token == "" {
		log.Write("not authenticated (no token specified)")
		return nil
	}

	user, res, err := s.client.Users.Get(ctx, "")
	if err != nil {
		if res == nil || res.StatusCode != http.StatusUnauthorized {
			return err
		}

		// TODO: rebuild client without token provider
		log.Write("not authenticated (token is invalid)")
		return nil
	}

	log.Write("authenticated as @%s", user.GetLogin())
	s.user = user

	if err := s.populateRepoCache(ctx, log); err != nil {
		return err
	}

	return nil
}

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func (s *source) populateRepoCache(
	ctx context.Context,
	log logs.Log,
) error {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	s.reposByID = map[int64]*github.Repository{}
	s.reposByOwner = map[string]map[string]*github.Repository{}

	for opts.Page != 0 {
		repoPage, res, err := s.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, r := range repoPage {
			log.WriteVerbose("discovered %s", r.GetFullName())

			owner := r.GetOwner().GetLogin()
			reposByName := s.reposByOwner[owner]
			if reposByName == nil {
				reposByName = map[string]*github.Repository{}
				s.reposByOwner[owner] = reposByName
			}

			reposByName[r.GetName()] = r
			s.reposByID[r.GetID()] = r
		}

		opts.Page = res.NextPage
	}

	log.Write(
		"added %d repositories to the repository list for @%s",
		len(s.reposByID),
		s.user.GetLogin(),
	)

	return nil
}
