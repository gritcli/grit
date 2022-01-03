package config

import "errors"

// GitHubConfig contains configuration specific to a GitHub repository source.
type GitHubConfig struct {
	// Domain is the base domain name of the GitHub installation.
	Domain string `hcl:"domain,optional"`

	// Token is a personal access token used to authenticate with the GitHub
	// API.
	Token string `hcl:"token,optional"`
}

// acceptVisitor calls v.VisitGitHubSource(s, c).
func (c GitHubConfig) acceptVisitor(s Source, v SourceVisitor) {
	v.VisitGitHubSource(s, c)
}

// withDefaults returns a copy of the configuration with any missing values
// replaced by their defaults.
func (c GitHubConfig) withDefaults() SourceConfig {
	if c.Domain == "" {
		c.Domain = "github.com"
	}

	return c
}

// validate returns an error if the configuration is invalid, it is intended
// to be called after any default values have been populated.
func (c GitHubConfig) validate() error {
	if c.Domain == "" {
		return errors.New("github domain must not be empty")
	}

	return nil
}

func init() {
	registerSourceType(
		"github",
		GitHubConfig{},
		Source{
			Name: "github",
			Config: GitHubConfig{
				Domain: "github.com",
			},
		},
	)
}
