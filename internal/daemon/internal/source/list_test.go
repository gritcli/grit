package source_test

import (
	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type List", func() {
	var list List

	Describe("func NewList()", func() {
		It("constructs the source from the configuration", func() {
			src := &stubs.Source{}

			list = NewList([]config.Source{
				{
					Name:    "<source>",
					Enabled: true,
					Clones: config.Clones{
						Dir: "/path/to/clones",
					},
					Driver: &stubs.SourceConfig{
						NewSourceFunc: func() sourcedriver.Source {
							return src
						},
					},
				},
			})

			Expect(list).To(Equal(List{
				{
					Name:        "<source>",
					Description: "<description>",
					CloneDir:    "/path/to/clones",
					Driver:      src,
				},
			}))

		})

		It("sorts sources by name", func() {
			list = NewList([]config.Source{
				{
					Name:    "<b>",
					Enabled: true,
					Driver:  &stubs.SourceConfig{},
				},
				{
					Name:    "<c>",
					Enabled: true,
					Driver:  &stubs.SourceConfig{},
				},
				{
					Name:    "<a>",
					Enabled: true,
					Driver:  &stubs.SourceConfig{},
				},
			})

			Expect(list).To(HaveLen(3))
			Expect(list[0].Name).To(Equal("<a>"))
			Expect(list[1].Name).To(Equal("<b>"))
			Expect(list[2].Name).To(Equal("<c>"))
		})

		It("excludes disabled sources", func() {
			list = NewList([]config.Source{
				{
					Name:    "<source>",
					Enabled: false,
				},
			})

			Expect(list).To(BeEmpty())
		})
	})

	Describe("func ByName()", func() {
		BeforeEach(func() {
			list = List{
				{
					Name: "<source>",
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
