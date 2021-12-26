package config

// GitHubDriver is a driver that supports GitHub.com and GitHub Enterprise
// Server installations.
const GitHubDriver SourceDriver = "github"

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
		GitHubDriver,
		GitHubConfig{},
		Source{
			Name:   "github.com",
			Driver: GitHubDriver,
			Config: GitHubConfig{
				Domain: "github.com",
			},
		},
	)
}
