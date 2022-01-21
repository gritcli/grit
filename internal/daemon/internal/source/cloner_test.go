package source_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Cloner", func() {
	var (
		tempDir string
		drv     *stubs.SourceDriver
		cloner  *Cloner
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())

		drv = &stubs.SourceDriver{}

		cloner = &Cloner{
			Sources: List{
				{
					Name:     "<source>",
					CloneDir: tempDir,
					Driver:   drv,
				},
			},
			Logger: logging.SilentLogger,
		}
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("func Clone()", func() {
		It("returns the clone directory", func() {
			drv.NewClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (sourcedriver.Cloner, string, error) {
				return &stubs.SourceDriverCloner{
					CloneFunc: func(
						context.Context,
						string,
						logging.Logger,
					) error {
						return nil
					},
				}, "clone-dir", nil
			}

			dir, err := cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dir).To(Equal(filepath.Join(tempDir, "clone-dir")))
		})

		It("returns an error if the directory already exists", func() {
			dir := filepath.Join(tempDir, "existing-dir")
			err := os.Mkdir(dir, 0700)
			Expect(err).ShouldNot(HaveOccurred())

			drv.NewClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (sourcedriver.Cloner, string, error) {
				return &stubs.SourceDriverCloner{}, "existing-dir", nil
			}

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
			drv.NewClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (sourcedriver.Cloner, string, error) {
				return &stubs.SourceDriverCloner{}, "\x00", nil
			}

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
			drv.NewClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (sourcedriver.Cloner, string, error) {
				return nil, "", errors.New("<error>")
			}

			_, err := cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).To(MatchError("unable to prepare for cloning: <error>"))
		})

		It("returns an error if the bound cloner returns an error", func() {
			drv.NewClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (sourcedriver.Cloner, string, error) {
				return &stubs.SourceDriverCloner{
					CloneFunc: func(
						context.Context,
						string,
						logging.Logger,
					) error {
						return errors.New("<error>")
					},
				}, "clone-dir", nil
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
