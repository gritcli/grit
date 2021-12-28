package githubdriver

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/config"
	"golang.org/x/oauth2"
)

// Source is an implementation of source.Source that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type Source struct {
	name   string
	domain string
	client *github.Client
	logger logging.Logger

	user  *github.User
	repos map[string]map[string]*github.Repository
}

// NewSource returns a new source with the given configuration.
func NewSource(
	name string,
	cfg config.GitHubConfig,
	logger logging.Logger,
) (*Source, error) {
	src := &Source{
		name:   name,
		domain: cfg.Domain,
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

	if isGitHubDotCom(cfg.Domain) {
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
func (s *Source) Name() string {
	return s.name
}

// Description returns a brief description of the repository source.
func (s *Source) Description() string {
	if isGitHubDotCom(s.domain) {
		return s.domain
	}

	return fmt.Sprintf("%s (github enterprise)", s.domain)
}

// Init initializes the source.
func (s *Source) Init(ctx context.Context) error {
	user, res, err := s.client.Users.Get(ctx, "")
	if err != nil {
		if res.StatusCode != http.StatusUnauthorized {
			return err
		}

		logging.Log(s.logger, "not authenticated")
		return nil
	}

	logging.Log(s.logger, "authenticated as %s", user.GetLogin())

	s.user = user

	if err := s.fetchRepos(ctx); err != nil {
		return err
	}

	return nil
}

// Run runs any background processes required by the source until ctx is
// canceled or a fatal error occurs.
func (s *Source) Run(ctx context.Context) error {
	return nil
}

func (s *Source) fetchRepos(ctx context.Context) error {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	for opts.Page != 0 {
		repos, res, err := s.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			owner := repo.GetOwner()

			if s.repos == nil {
				s.repos = map[string]map[string]*github.Repository{}
			}

			repoMap := s.repos[owner.GetLogin()]
			if repoMap == nil {
				repoMap = map[string]*github.Repository{}
				s.repos[owner.GetLogin()] = repoMap
			}

			logging.Log(s.logger, "indexed repository: %s", repo.GetFullName())
			repoMap[repo.GetName()] = repo
		}

		opts.Page = res.NextPage
	}

	return nil
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(domain string) bool {
	return strings.EqualFold(domain, "github.com")
}
