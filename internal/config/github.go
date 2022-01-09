package config

import (
	"strings"
)

// GitHub contains configuration specific to a GitHub repository source.
type GitHub struct {
	// Domain is the base domain name of the GitHub installation.
	Domain string

	// Token is a personal access token used to authenticate with the GitHub
	// API.
	Token string

	// Git is the configuration that controls how Grit uses Git for this source.
	Git Git
}

// acceptVisitor calls v.VisitGitHubSource(s, c).
func (c GitHub) acceptVisitor(s Source, v SourceVisitor) {
	v.VisitGitHubSource(s, c)
}

// String returns a human-readable description of the configuration.
func (c GitHub) String() string {
	desc := c.Domain

	if !strings.EqualFold(c.Domain, "github.com") {
		desc += " (github enterprise server)"
	}

	return desc
}

// gitHubBlock is the HCL schema for a "source" block that uses the "github"
// source driver.
type gitHubBlock struct {
	Domain string    `hcl:"domain,optional"`
	Token  string    `hcl:"token,optional"`
	Git    *gitBlock `hcl:"git,block"`
}

func (b *gitHubBlock) Normalize(cfg unresolvedConfig, s unresolvedSource) error {
	if b.Domain == "" {
		b.Domain = "github.com"
	}

	return normalizeSourceSpecificGitBlock(cfg, s, &b.Git)
}

func (b *gitHubBlock) Assemble() SourceDriverConfig {
	return GitHub{
		Domain: b.Domain,
		Token:  b.Token,
		Git:    assembleGitBlock(*b.Git),
	}
}

func init() {
	registerSourceDriver(
		"github",
		func() sourceDriverBlock {
			return &gitHubBlock{}
		},
	)

	registerDefaultSource(
		"github",
		func() sourceDriverBlock {
			return &gitHubBlock{
				Domain: "github.com",
			}
		},
	)
}
