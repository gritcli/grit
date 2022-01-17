package source_test

import (
	"github.com/gritcli/grit/internal/daemon/internal/config"
	. "github.com/gritcli/grit/internal/daemon/internal/source"
	. "github.com/gritcli/grit/internal/daemon/internal/source/internal/fixtures"
	"github.com/gritcli/grit/plugin/driver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type List", func() {
	var list List

	Describe("func NewList()", func() {
		It("constructs the source from the configuration", func() {
			d := &DriverStub{}

			list = NewList([]config.Source{
				{
					Name:    "<source>",
					Enabled: true,
					Clones: config.Clones{
						Dir: "/path/to/clones",
					},
					Driver: &DriverConfigStub{
						NewDriverFunc: func() driver.Driver {
							return d
						},
					},
				},
			})

			Expect(list).To(Equal(List{
				{
					Name:        "<source>",
					Description: "<driver config stub>",
					CloneDir:    "/path/to/clones",
					Driver:      d,
				},
			}))

		})

		It("sorts sources by name", func() {
			list = NewList([]config.Source{
				{
					Name:    "<b>",
					Enabled: true,
					Driver:  &DriverConfigStub{},
				},
				{
					Name:    "<c>",
					Enabled: true,
					Driver:  &DriverConfigStub{},
				},
				{
					Name:    "<a>",
					Enabled: true,
					Driver:  &DriverConfigStub{},
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
