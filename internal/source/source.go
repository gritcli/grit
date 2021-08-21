package source

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v38/github"
	"github.com/gritcli/grit/internal/config"
	"github.com/gritcli/grit/internal/source/githubsource"
	"github.com/gritcli/grit/internal/source/vanillasource"
	"golang.org/x/oauth2"
)

// Source is an interface for a source of repositories.
type Source interface {
	// Description returns a short, human-readable description of the source.
	//
	// The description should be adequate to distinguish this source from any
	// other sources that may exist.
	Description() string

	// Status queries the status of the source.
	//
	// It returns an error if the source is misconfigured or unreachable.
	//
	// The status string should include any source-specific information
	Status(ctx context.Context) (string, error)
}

// FromConfig creates a new source from a source configuration element.
func FromConfig(src config.Source) Source {
	var f factory
	if err := src.Visit(&f); err != nil {
		panic(err)
	}

	return f.Result
}

// factory is an implementation of config.SourceVisitor that constructs sources
// from a config.Source element.
type factory struct {
	Result Source
}

func (f *factory) VisitGitSource(src config.GitSource) error {
	f.Result = &vanillasource.Source{
		Endpoint: src.Endpoint,
	}

	return nil
}
func (f *factory) VisitGitHubSource(src config.GitHubSource) error {
	hc := http.DefaultClient
	if src.Token != "" {
		hc = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: src.Token,
			}),
		)
	}

	u, err := url.Parse(src.API.String()) // clone URL
	if err != nil {
		return err
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	gc := github.NewClient(hc)
	gc.BaseURL = u

	f.Result = &githubsource.Source{
		Client: gc,
	}

	return nil
}
