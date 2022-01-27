package githubsource_test

import (
	"path/filepath"

	"github.com/gritcli/grit/driver/sourcedriver"
)

var (
	publicOrgRepo = sourcedriver.RemoteRepo{
		ID:               "451303002",
		Name:             "grit-integration-tests-org/test-public",
		Description:      "Used to test that Grit works with public GitHub repositories that belong to an organization.",
		WebURL:           "https://github.com/grit-integration-tests-org/test-public",
		RelativeCloneDir: filepath.Join("grit-integration-tests-org", "test-public"),
	}

	privateOrgRepo = sourcedriver.RemoteRepo{
		ID:               "451303236",
		Name:             "grit-integration-tests-org/test-private",
		Description:      "Used to test that Grit works with private GitHub repositories that belong to an organization.",
		WebURL:           "https://github.com/grit-integration-tests-org/test-private",
		RelativeCloneDir: filepath.Join("grit-integration-tests-org", "test-private"),
	}

	publicUserRepo = sourcedriver.RemoteRepo{
		ID:               "451288349",
		Name:             "grit-integration-tests/test-public",
		Description:      "Used to test that Grit works with public GitHub repositories.",
		WebURL:           "https://github.com/grit-integration-tests/test-public",
		RelativeCloneDir: filepath.Join("grit-integration-tests", "test-public"),
	}

	privateUserRepo = sourcedriver.RemoteRepo{
		ID:               "451288389",
		Name:             "grit-integration-tests/test-private",
		Description:      "Used to test that Grit works with private GitHub repositories.",
		WebURL:           "https://github.com/grit-integration-tests/test-private",
		RelativeCloneDir: filepath.Join("grit-integration-tests", "test-private"),
	}

	// thirdPartyRepo is a repository that the authenticated user does not have
	// access to.
	//
	// The CI process uses a GitHub personal access token belonging to
	// @grit-integration-tests which is NOT a member of the "grit-cli"
	// organization.
	thirdPartyRepo = sourcedriver.RemoteRepo{
		ID:               "397822937",
		Name:             "gritcli/grit",
		Description:      "Manage your local Git clones.",
		WebURL:           "https://github.com/gritcli/grit",
		RelativeCloneDir: filepath.Join("gritcli", "grit"),
	}

	// allTestRepos is the set of all test repositories to which the
	// @grit-integration-tests GitHub user has been granted explicit access.
	allTestRepos = []sourcedriver.RemoteRepo{
		publicOrgRepo,
		privateOrgRepo,
		publicUserRepo,
		privateUserRepo,
	}
)
