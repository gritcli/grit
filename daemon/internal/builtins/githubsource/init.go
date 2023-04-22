package githubsource

import (
	"context"
	"net/http"

	"github.com/google/go-github/v50/github"
	"github.com/gritcli/grit/daemon/internal/logs"
	"golang.org/x/oauth2"
)

// Init initializes the source.
func (s *source) Init(
	ctx context.Context,
	log logs.Log,
) error {
	return s.init(
		ctx,
		s.config.Token,
		log,
	)
}

func (s *source) init(
	ctx context.Context,
	token string,
	log logs.Log,
) error {
	httpClient := http.DefaultClient
	if token != "" {
		httpClient = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			),
		)
	}

	client, err := newClient(s.config, httpClient)
	if err != nil {
		return err
	}

	state := &state{
		Client: client,
	}
	s.state.Store(state)

	if s.config.Token == "" {
		log.Write("not authenticated (no token specified)")
		return nil
	}

	user, res, err := state.Client.Users.Get(ctx, "")
	if err != nil {
		if res == nil || res.StatusCode != http.StatusUnauthorized {
			return err
		}

		// TODO: rebuild client without token provider
		log.Write("not authenticated (token is invalid)")
		return nil
	}

	log.Write("authenticated as @%s", user.GetLogin())
	state.User = user

	if err := populateRepoCache(ctx, state, log); err != nil {
		return err
	}

	return nil
}

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func populateRepoCache(
	ctx context.Context,
	state *state,
	log logs.Log,
) error {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	state.ReposByID = map[int64]*github.Repository{}
	state.ReposByOwner = map[string]map[string]*github.Repository{}

	for opts.Page != 0 {
		repoPage, res, err := state.Client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, r := range repoPage {
			log.WriteVerbose("discovered %s", r.GetFullName())

			owner := r.GetOwner().GetLogin()
			reposByName := state.ReposByOwner[owner]
			if reposByName == nil {
				reposByName = map[string]*github.Repository{}
				state.ReposByOwner[owner] = reposByName
			}

			reposByName[r.GetName()] = r
			state.ReposByID[r.GetID()] = r
		}

		opts.Page = res.NextPage
	}

	log.Write(
		"added %d repositories to the repository list for @%s",
		len(state.ReposByID),
		state.User.GetLogin(),
	)

	return nil
}
