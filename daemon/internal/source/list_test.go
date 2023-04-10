package source_test

import (
	"context"
	"net/url"

	"github.com/gritcli/grit/daemon/internal/config"
	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
	. "github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/daemon/internal/stubs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type List", func() {
	var list List

	Describe("func NewList()", func() {
		It("constructs the source from the configuration", func() {
			// These sources have an InitFunc set so that they can be
			// differentiated from each other. The panic message is different so
			// these functions can't get combined by the compiler.
			srcA := &stubs.Source{
				InitFunc: func(context.Context, sourcedriver.InitParameters, logs.Log) error {
					panic("not implemented (a)")
				},
			}

			srcB := &stubs.Source{
				InitFunc: func(context.Context, sourcedriver.InitParameters, logs.Log) error {
					panic("not implemented (b)")
				},
			}

			baseURL, err := url.Parse("http://localhost:8080")
			Expect(err).ShouldNot(HaveOccurred())

			list = NewList(
				baseURL,
				[]config.Source{
					{
						Name:    "<source-a>",
						Enabled: true,
						Clones: config.Clones{
							Dir: "/path/to/clones-a",
						},
						Driver: &stubs.SourceConfig{
							NewSourceFunc: func() sourcedriver.Source {
								return srcA
							},
						},
					},
					{
						Name:    "<source-b>",
						Enabled: true,
						Clones: config.Clones{
							Dir: "/path/to/clones-b",
						},
						Driver: &stubs.SourceConfig{
							NewSourceFunc: func() sourcedriver.Source {
								return srcB
							},
						},
					},
				},
			)

			Expect(list).To(ConsistOf(
				Source{
					Name:         "<source-a>",
					Description:  "<description>",
					BaseCloneDir: "/path/to/clones-a",
					BaseURL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "source/<source-a>",
					},
					Driver: srcA,
				},
				Source{
					Name:         "<source-b>",
					Description:  "<description>",
					BaseCloneDir: "/path/to/clones-b",
					BaseURL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "source/<source-b>",
					},
					Driver: srcB,
				},
			))

		})

		It("excludes disabled sources", func() {
			baseURL, err := url.Parse("http://localhost:8080")
			Expect(err).ShouldNot(HaveOccurred())

			list = NewList(
				baseURL,
				[]config.Source{
					{
						Name:    "<source>",
						Enabled: false,
					},
				},
			)

			Expect(list).To(BeEmpty())
		})
	})

	Describe("func ByName()", func() {
		BeforeEach(func() {
			list = List{
				{
					Name: "<source>",
				},
				{
					Name: "<other>",
				},
			}
		})

		It("returns the source with the given name (case-insensitive)", func() {
			s, ok := list.ByName("<SOURCE>")
			Expect(ok).To(BeTrue())
			Expect(s).To(Equal(Source{
				Name: "<source>",
			}))
		})

		It("returns false if the source is not in the list", func() {
			_, ok := list.ByName("<unknown>")
			Expect(ok).To(BeFalse())
		})
	})
})
