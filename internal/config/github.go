package config

// GitHubDriver is a driver that supports GitHub.com and GitHub Enterprise
// Server installations.
const GitHubDriver SourceDriver = "github"

// GitHubConfig contains configuration specific to a GitHub repository source.
type GitHubConfig struct {
	// Domain is the base domain name of the GitHub installation.
	Domain string `hcl:"domain,optional"`
}

// DescribeConfig returns a short, human-readable description of the
// configuration.
func (c GitHubConfig) DescribeConfig() string {
	return c.Domain
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
