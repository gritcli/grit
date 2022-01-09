package github_test

import (
	"context"
	"math"
	"os"
	"strconv"

	"github.com/dogmatiq/dodeca/logging"
	git "github.com/go-git/go-git/v5"
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func source.Clone()", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		configure func(*config.GitHub)
		driver    source.Driver
		tempDir   string
		logger    logging.DiscardLogger
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	When("unauthenticated", func() {
		BeforeEach(func() {
			configure = func(*config.GitHub) {}
		})

		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachUnauthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		When("HTTP is the preferred protocol", func() {
			BeforeEach(func() {
				configure = func(cfg *config.GitHub) {
					cfg.Git.PreferHTTP = true
				}
			})

			It("clones the repository using HTTP", func() {
				cloneDir, err := driver.Clone(ctx, gritRepo.ID, tempDir, logger)
				skipIfRateLimited(err)
				Expect(cloneDir).To(Equal(gritRepo.Name))

				repo, err := git.PlainOpen(tempDir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("https://github.com/gritcli/grit.git"))
			})
		})

		When("the SSH agent is unavailable", func() {
			var orig string
			BeforeEach(func() {
				orig = os.Getenv("SSH_AUTH_SOCK")
				os.Setenv("SSH_AUTH_SOCK", "")
			})

			AfterEach(func() {
				os.Setenv("SSH_AUTH_SOCK", orig)
			})

			It("clones the repository using HTTP", func() {
				cloneDir, err := driver.Clone(ctx, gritRepo.ID, tempDir, logger)
				skipIfRateLimited(err)
				Expect(cloneDir).To(Equal(gritRepo.Name))

				repo, err := git.PlainOpen(tempDir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("https://github.com/gritcli/grit.git"))
			})
		})

		When("using an explicit private key", func() {
			BeforeEach(func() {
				configure = func(cfg *config.GitHub) {
					cfg.Git.SSHKeyFile = "./testdata/deploy-key-no-passphrase"
				}
			})

			It("clones the repository using SSH", func() {
				cloneDir, err := driver.Clone(ctx, gritRepo.ID, tempDir, logger)
				skipIfRateLimited(err)
				Expect(cloneDir).To(Equal(gritRepo.Name))

				repo, err := git.PlainOpen(tempDir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("git@github.com:gritcli/grit.git"))
			})

			When("the private key file can not be laoded", func() {
				BeforeEach(func() {
					configure = func(cfg *config.GitHub) {
						cfg.Git.SSHKeyFile = "/does/not/exist"
					}
				})

				It("returns an error", func() {
					_, err := driver.Clone(ctx, gritRepo.ID, tempDir, logger)
					Expect(err).To(MatchError("open /does/not/exist: no such file or directory"))
				})
			})
		})

		When("using an explicit private key with a passphrase", func() {
			BeforeEach(func() {
				configure = func(cfg *config.GitHub) {
					cfg.Git.SSHKeyFile = "./testdata/deploy-key-with-passphrase"
					cfg.Git.SSHKeyPassphrase = "passphrase"
				}
			})

			It("clones the repository using SSH", func() {
				cloneDir, err := driver.Clone(ctx, gritRepo.ID, tempDir, logger)
				skipIfRateLimited(err)
				Expect(cloneDir).To(Equal(gritRepo.Name))

				repo, err := git.PlainOpen(tempDir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("git@github.com:gritcli/grit.git"))
			})
		})

		When("the SSH agent is available", func() {
			BeforeEach(func() {
				if os.Getenv("SSH_AUTH_SOCK") == "" {
					Skip("the SSH agent is not available")
				}
			})

			It("clones the repository using SSH", func() {
				cloneDir, err := driver.Clone(ctx, gritRepo.ID, tempDir, logger)
				skipIfRateLimited(err)
				Expect(cloneDir).To(Equal(gritRepo.Name))

				repo, err := git.PlainOpen(tempDir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("git@github.com:gritcli/grit.git"))
			})
		})

		It("returns an error if the repository does not exist", func() {
			_, err := driver.Clone(ctx, strconv.FormatInt(math.MaxInt64, 10), tempDir, logger)
			Expect(err).To(MatchError("GET https://api.github.com/repositories/9223372036854775807: 404 Not Found []"))
		})

		It("returns an error if the repository ID is invalid", func() {
			_, err := driver.Clone(ctx, "<invalid>", tempDir, logger)
			Expect(err).To(MatchError("invalid repo ID, expected positive integer"))
		})

		It("returns an error if the repository ID is non-positive", func() {
			_, err := driver.Clone(ctx, "0", tempDir, logger)
			Expect(err).To(MatchError("invalid repo ID, expected positive integer"))
		})
	})

	When("authenticated", func() {
		BeforeEach(func() {
			configure = func(*config.GitHub) {}
		})

		JustBeforeEach(func() {
			ctx, cancel, driver = beforeEachAuthenticated(configure)
		})

		AfterEach(func() {
			cancel()
		})

		When("HTTP is the preferred protocol", func() {
			BeforeEach(func() {
				configure = func(cfg *config.GitHub) {
					cfg.Git.PreferHTTP = true
				}
			})

			It("clones the repository using HTTP with token-based authentication", func() {
				// TODO: https://github.com/gritcli/grit/issues/13
				//
				// Change this test to use a private repository so it's actually
				// verifying that the token is being used.

				cloneDir, err := driver.Clone(ctx, gritRepo.ID, tempDir, logger)
				skipIfRateLimited(err)
				Expect(cloneDir).To(Equal(gritRepo.Name))

				repo, err := git.PlainOpen(tempDir)
				Expect(err).ShouldNot(HaveOccurred())

				rem, err := repo.Remote("origin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rem.Config().URLs).To(ConsistOf("https://github.com/gritcli/grit.git"))
			})
		})
	})
})
