package source_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dogmatiq/dodeca/logging"
	. "github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/internal/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Cloner", func() {
	var (
		tempDir      string
		repo         sourcedriver.RemoteRepo
		sourceCloner *stubs.SourceCloner
		driver       *stubs.Source
		src          Source
		cloner       *Cloner
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())

		repo = sourcedriver.RemoteRepo{
			ID:               "<id>",
			Name:             "<repo>",
			RelativeCloneDir: "clone-dir",
		}

		sourceCloner = &stubs.SourceCloner{}

		driver = &stubs.Source{
			NewClonerFunc: func(
				context.Context,
				string,
				logging.Logger,
			) (sourcedriver.Cloner, sourcedriver.RemoteRepo, error) {
				return sourceCloner, repo, nil
			},
		}

		src = Source{
			Name:     "<source>",
			CloneDir: tempDir,
			Driver:   driver,
		}

		cloner = &Cloner{
			Sources: List{src},
			Logger:  logging.SilentLogger,
		}
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("func Clone()", func() {
		It("returns the local repo", func() {
			local, err := cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(local).To(Equal(
				LocalRepo{
					RemoteRepo:       repo,
					Source:           src,
					AbsoluteCloneDir: filepath.Join(tempDir, "clone-dir"),
				},
			))
		})

		It("returns an error if the directory already exists", func() {
			dir := filepath.Join(tempDir, "clone-dir")
			err := os.Mkdir(dir, 0700)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).To(MatchError(
				fmt.Sprintf(
					"unable to create clone directory: mkdir %s: file exists",
					dir,
				),
			))
		})

		It("returns an error if the directory can not be created", func() {
			repo.RelativeCloneDir = "\x00"

			_, err := cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).To(MatchError(
				fmt.Sprintf(
					"unable to create clone directory: stat %s/\x00: invalid argument",
					tempDir,
				),
			))
		})

		It("returns an error if the source is not known", func() {
			_, err := cloner.Clone(
				context.Background(),
				"<unknown>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).To(MatchError("unable to clone: unrecognized source (<unknown>)"))
		})

		It("returns an error if the driver returns an error", func() {
			driver.NewClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (sourcedriver.Cloner, sourcedriver.RemoteRepo, error) {
				return nil, sourcedriver.RemoteRepo{}, errors.New("<error>")
			}

			_, err := cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).To(MatchError("unable to prepare for cloning: <error>"))
		})

		It("returns an error if the source's cloner returns an error", func() {
			sourceCloner.CloneFunc = func(
				context.Context,
				string,
				logging.Logger,
			) error {
				return errors.New("<error>")
			}

			_, err := cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).To(MatchError("unable to clone: <error>"))

			_, err = os.Stat(
				filepath.Join(tempDir, "clone-dir"),
			)
			Expect(err).Should(HaveOccurred())
			Expect(os.IsNotExist(err)).To(BeTrue(), err.Error())
		})
	})
})
