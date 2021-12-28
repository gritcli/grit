package config

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

// withDefaults returns a copy of the configuration populated with default
// values.
func (c GitHubConfig) withDefaults() DriverConfig {
	if c.Domain == "" {
		c.Domain = "github.com"
	}

	return c
}

func init() {
	registerDriver(
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
