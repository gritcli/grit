package config

// GitHubConfig contains configuration specific to a GitHub repository source.
type GitHubConfig struct {
	// Domain is the base domain name of the GitHub installation.
	Domain string

	// Token is a personal access token used to authenticate with the GitHub
	// API.
	Token string
}

// acceptVisitor calls v.VisitGitHubSource(s, c).
func (c GitHubConfig) acceptVisitor(s Source, v SourceVisitor) {
	v.VisitGitHubSource(s, c)
}

// gitHubBlock is the HCL schema for a "source" block that uses the "github"
// source implementation.
type gitHubBlock struct {
	Domain string `hcl:"domain,optional"`
	Token  string `hcl:"token,optional"`
}

func (b gitHubBlock) resolve(filename string, cfg Config) (SourceConfig, error) {
	c := GitHubConfig(b)

	if c.Domain == "" {
		c.Domain = "github.com"
	}

	return c, nil
}

func init() {
	registerSourceSchema(
		"github",
		gitHubBlock{},
		Source{
			Name:    "github",
			Enabled: true,
			Config: GitHubConfig{
				Domain: "github.com",
			},
		},
	)
}
