package config

// GitHubConfig contains configuration specific to a GitHub repository source.
type GitHubConfig struct {
	// Domain is the base domain name of the GitHub installation.
	Domain string `hcl:"domain,optional"`
}

// String returns a short, human-readable description of the configuration.
func (c GitHubConfig) String() string {
	return c.Domain
}

// acceptVisitor calls v.VisitGitHubSource(s, c).
func (c GitHubConfig) acceptVisitor(s Source, v SourceVisitor) {
	v.VisitGitHubSource(s, c)
}

func init() {
	registerDriver(
		"github",
		GitHubConfig{},
		Source{
			Name: "github.com",
			Config: GitHubConfig{
				Domain: "github.com",
			},
		},
	)
}
