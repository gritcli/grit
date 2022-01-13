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
		drv       *Driver
	)

	BeforeEach(func() {
		configure = func(*config.GitHub) {}
	})

	AfterEach(func() {
		cancel()
	})

	When("unauthenticated", func() {
		JustBeforeEach(func() {
			ctx, cancel, drv = beforeEachUnauthenticated(configure)
		})

		It("returns a gitvcs.Cloner", func() {
			cloner, dir, err := drv.NewCloner(ctx, gritPublicTestRepo.ID, logging.SilentLogger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:      "git@github.com:gritcli/test-public.git",
				SSHKeyFile:       drv.Config.Git.SSHKeyFile,
				SSHKeyPassphrase: drv.Config.Git.SSHKeyPassphrase,
				HTTPEndpoint:     "https://github.com/gritcli/test-public.git",
				PreferHTTP:       drv.Config.Git.PreferHTTP,
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})

	When("authenticated", func() {
		JustBeforeEach(func() {
			ctx, cancel, drv = beforeEachAuthenticated(configure)
		})

		It("returns a gitvcs.Cloner with the token as the HTTP password", func() {
			// TODO: https://github.com/gritcli/grit/issues/13
			//
			// Test with a private repository instead.
			cloner, dir, err := drv.NewCloner(ctx, gritPublicTestRepo.ID, logging.SilentLogger)
			skipIfRateLimited(err)

			Expect(cloner).To(Equal(&gitvcs.Cloner{
				SSHEndpoint:      "git@github.com:gritcli/test-public.git",
				SSHKeyFile:       drv.Config.Git.SSHKeyFile,
				SSHKeyPassphrase: drv.Config.Git.SSHKeyPassphrase,
				HTTPEndpoint:     "https://github.com/gritcli/test-public.git",
				HTTPPassword:     drv.Config.Token,
				PreferHTTP:       drv.Config.Git.PreferHTTP,
			}))

			Expect(dir).To(Equal("gritcli/test-public"))
		})
	})
})
