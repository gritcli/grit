package config

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

// gitHubBlock is the HCL schema for a "source" block that uses the "github"
// source implementation.
type gitHubBlock struct {
	Domain string    `hcl:"domain,optional"`
	Token  string    `hcl:"token,optional"`
	Git    *gitBlock `hcl:"git,block"`
}

func (b *gitHubBlock) Normalize(cfg unresolvedConfig) error {
	if b.Domain == "" {
		b.Domain = "github.com"
	}

	if b.Git == nil {
		b.Git = &gitBlock{}
	}

	if b.Git.PrivateKey == "" {
		b.Git.PrivateKey = cfg.GlobalGit.Block.PrivateKey
	}

	if b.Git.PreferHTTP == nil {
		b.Git.PreferHTTP = cfg.GlobalGit.Block.PreferHTTP
	}

	return nil
}

func (b *gitHubBlock) Assemble() SourceConfig {
	return GitHub{
		Domain: b.Domain,
		Token:  b.Token,
		Git:    assembleGitBlock(*b.Git),
	}
}

func init() {
	registerSourceImpl(
		"github",
		func() sourceBlockBody {
			return &gitHubBlock{}
		},
	)

	registerDefaultSource(
		"github",
		func() sourceBlockBody {
			return &gitHubBlock{
				Domain: "github.com",
			}
		},
	)
}
