package github_test

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/daemon/internal/config"
	. "github.com/gritcli/grit/internal/daemon/internal/source/internal/github"
	"github.com/gritcli/grit/plugin/vcs/gitvcs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Driver.NewCloner()", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		configure func(*config.GitHub)
		driver    *Driver
	)

	BeforeEach(func() {
		configure = func(*config.GitHub) {}
	})

	AfterEach(func() {
		cancel()
	})

	When("unauthenticated", func() {
		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachUnauthenticated(configure)
		})

		It("returns a gitvcs.Cloner", func() {
			cloner, dir, err := driver.NewCloner(ctx, gritPublicTestRepo.ID, logging.SilentLogger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:      "git@github.com:gritcli/test-public.git",
				SSHKeyFile:       driver.Config.Git.SSHKeyFile,
				SSHKeyPassphrase: driver.Config.Git.SSHKeyPassphrase,
				HTTPEndpoint:     "https://github.com/gritcli/test-public.git",
				PreferHTTP:       driver.Config.Git.PreferHTTP,
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})

	When("authenticated", func() {
		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachAuthenticated(configure)
		})

		It("returns a git.BoundCloner with the token as the HTTP password", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// Test with a private repository instead.
			cloner, dir, err := driver.NewCloner(ctx, gritPublicTestRepo.ID, logging.SilentLogger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:      "git@github.com:gritcli/test-public.git",
				SSHKeyFile:       driver.Config.Git.SSHKeyFile,
				SSHKeyPassphrase: driver.Config.Git.SSHKeyPassphrase,
				HTTPEndpoint:     "https://github.com/gritcli/test-public.git",
				HTTPPassword:     driver.Config.Token,
				PreferHTTP:       driver.Config.Git.PreferHTTP,
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})
})
