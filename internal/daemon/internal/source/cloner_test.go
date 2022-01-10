package source_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/dogmatiq/dodeca/logging"
	. "github.com/gritcli/grit/internal/daemon/internal/source"
	. "github.com/gritcli/grit/internal/daemon/internal/source/internal/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Cloner", func() {
	var (
		driver *DriverStub
		cloner *Cloner
	)

	BeforeEach(func() {
		driver = &DriverStub{}

		cloner = &Cloner{
			Sources: List{
				{
					Name:   "<source>",
					Driver: driver,
				},
			},
			Logger: logging.SilentLogger,
		}
	})

	Describe("func Clone()", func() {
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
			driver.NewBoundClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (BoundCloner, string, error) {
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
			dir, err := os.MkdirTemp("", "")
			Expect(err).ShouldNot(HaveOccurred())
			defer os.RemoveAll(dir)

			dir = filepath.Join(dir, "clone-dir")

			driver.NewBoundClonerFunc = func(
				context.Context,
				string,
				logging.Logger,
			) (BoundCloner, string, error) {
				return &BoundClonerStub{
					CloneFunc: func(c context.Context, s string) error {
						return errors.New("<error>")
					},
				}, dir, nil
			}

			_, err = cloner.Clone(
				context.Background(),
				"<source>",
				"<id>",
				logging.SilentLogger,
			)
			Expect(err).To(MatchError("unable to clone: <error>"))
		})
	})
})
