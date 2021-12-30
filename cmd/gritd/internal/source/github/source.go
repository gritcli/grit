package github

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

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
	repoCache  map[string]map[string]source.Repo
}

// NewSource returns a new source with the given configuration.
func NewSource(
	name string,
	cfg config.GitHubConfig,
	logger logging.Logger,
) (source.Source, error) {
	src := &impl{
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
func (s *impl) Name() string {
	return s.name
}

// Description returns a brief description of the repository source.
func (s *impl) Description() string {
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

	s.userM.Lock()
	s.user = user
	s.userM.Unlock()

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

	repos := map[string]map[string]source.Repo{}
	count := 0

	for opts.Page != 0 {
		repoPage, res, err := s.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return err
		}

		for _, r := range repoPage {
			owner := r.GetOwner()

			reposByOwner := repos[owner.GetLogin()]
			if reposByOwner == nil {
				reposByOwner = map[string]source.Repo{}
				repos[owner.GetLogin()] = reposByOwner
			}

			logging.Debug(s.logger, "cached repository: %s", r.GetFullName())
			count++

			reposByOwner[r.GetName()] = convertRepo(r)
		}

		opts.Page = res.NextPage
	}

	s.repoCacheM.Lock()
	s.repoCache = repos
	s.repoCacheM.Unlock()

	logging.Log(s.logger, "cached %d repositories across %d owner(s)", count, len(repos))

	return nil
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(domain string) bool {
	return strings.EqualFold(domain, "github.com")
}

// convertRepo converts a github.Repository to a source.Repo.
func convertRepo(r *github.Repository) source.Repo {
	return source.Repo{
		ID:          strconv.FormatInt(r.GetID(), 10),
		Name:        r.GetFullName(),
		Description: r.GetDescription(),
		WebURL:      r.GetHTMLURL(),
	}
}