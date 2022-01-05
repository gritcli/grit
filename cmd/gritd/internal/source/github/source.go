package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/config"
	"golang.org/x/oauth2"
)

// impl is an implementation of source.Source that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type impl struct {
	name   string
	cfg    config.GitHub
	client *github.Client
	logger logging.Logger

	cache cache
}

// NewSource returns a new source with the given configuration.
func NewSource(
	name string,
	cfg config.GitHub,
	logger logging.Logger,
) (source.Source, error) {
	src := &impl{
		name:   name,
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
		src.client = github.NewClient(httpClient)
	} else {
		var err error
		src.client, err = github.NewEnterpriseClient(cfg.Domain, "", httpClient)
		if err != nil {
			return nil, err
		}
	}

	return src, nil
}

// Name returns a short, human-readable identifier of the repository source.
func (s *impl) Name() string {
	return s.name
}

// Description returns a brief description of the repository source.
func (s *impl) Description() string {
	var info []string

	if !isGitHubDotCom(s.cfg) {
		info = append(info, "github enterprise")
	}

	if u := s.cache.CurrentUser(); u != nil {
		info = append(info, "@"+u.GetLogin())
	}

	if len(info) > 0 {
		return fmt.Sprintf("%s (%s)", s.cfg.Domain, strings.Join(info, ", "))
	}

	return s.cfg.Domain
}

// Init initializes the source.
func (s *impl) Init(ctx context.Context) error {
	user, res, err := s.client.Users.Get(ctx, "")
	if err != nil {
		if res.StatusCode != http.StatusUnauthorized {
			return err
		}

		logging.Log(s.logger, "not authenticated")
		return nil
	}

	logging.Log(s.logger, "authenticated as %s", user.GetLogin())
	s.cache.SetCurrentUser(user)

	if err := s.populateRepoCache(ctx); err != nil {
		return err
	}

	return nil
}

// Run performs any background tasks required by the source.
func (s *impl) Run(ctx context.Context) error {
	return nil
}

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func (s *impl) populateRepoCache(ctx context.Context) error {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var repos []*github.Repository

	for opts.Page != 0 {
		repoPage, res, err := s.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, r := range repoPage {
			logging.Debug(s.logger, "cached repository: %s", r.GetFullName())
			repos = append(repos, r)
		}

		opts.Page = res.NextPage
	}

	logging.Log(
		s.logger,
		"cached %d repositories",
		len(repos),
	)

	s.cache.SetRepos(repos)

	return nil
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(cfg config.GitHub) bool {
	return strings.EqualFold(cfg.Domain, "github.com")
}
