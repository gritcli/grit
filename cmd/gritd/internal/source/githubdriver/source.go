package githubdriver

import (
	"context"
	"fmt"
	"strings"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/google/go-github/github"
	"github.com/gritcli/grit/internal/config"
)

// Source is an implementation of source.Source that provides repositories from
// GitHub.com or a GitHub Enterprise Server installation.
type Source struct {
	name   string
	domain string
	client *github.Client
	logger logging.Logger
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

	if isGitHubDotCom(cfg.Domain) {
		src.client = github.NewClient(nil)
	} else {
		var err error
		src.client, err = github.NewEnterpriseClient(cfg.Domain, "", nil)
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

// Run runs any background processes required by the source until ctx is
// canceled or a fatal error occurs.
func (s *Source) Run(ctx context.Context) error {
	logging.Log(s.logger, s.Description())
	return nil
}

// isGitHubDotCom returns true if domain is the domain for github.com.
func isGitHubDotCom(domain string) bool {
	return strings.EqualFold(domain, "github.com")
}
