package githubdriver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

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

	// user is the authenticated user, based on the token in the source
	// configuration. If no token is provided, or it is invalid, user is nil.
	userM sync.RWMutex
	user  *github.User

	// repoCache is an in-memory cache of the repositores to which the
	// authenticated user has explicit read, wrote or admin access.
	repoCacheM sync.RWMutex
	repoCache  map[string]map[string]*github.Repository
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
	var info []string

	if !isGitHubDotCom(s.domain) {
		info = append(info, "github enterprise")
	}

	s.userM.RLock()
	user := s.user
	s.userM.RUnlock()

	if user != nil {
		info = append(info, "@"+user.GetLogin())
	}

	if len(info) > 0 {
		return fmt.Sprintf("%s (%s)", s.domain, strings.Join(info, ", "))
	}

	return s.domain
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

	s.userM.Lock()
	s.user = user
	s.userM.Unlock()

	if err := s.populateRepoCache(ctx); err != nil {
		return err
	}

	return nil
}

// Run runs any background processes required by the source until ctx is
// canceled or a fatal error occurs.
func (s *Source) Run(ctx context.Context) error {
	return nil
}

// populateRepoCache populates s.populateRepoCache with the repositories to
// which the authenticated user has explicit read, write or admin access.
func (s *Source) populateRepoCache(ctx context.Context) error {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	repos := map[string]map[string]*github.Repository{}

	for opts.Page != 0 {
		repoPage, res, err := s.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, r := range repoPage {
			owner := r.GetOwner()

			reposByOwner := repos[owner.GetLogin()]
			if reposByOwner == nil {
				reposByOwner = map[string]*github.Repository{}
				repos[owner.GetLogin()] = reposByOwner
			}

			logging.Log(s.logger, "cached repository: %s", r.GetFullName())
			reposByOwner[r.GetName()] = r
		}

		opts.Page = res.NextPage
	}

	s.repoCacheM.Lock()
	s.repoCache = repos
	s.repoCacheM.Unlock()

	return nil
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(domain string) bool {
	return strings.EqualFold(domain, "github.com")
}
