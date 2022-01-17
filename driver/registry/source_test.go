package registry_test

import (
	. "github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Registry", func() {
	var reg *Registry

	BeforeEach(func() {
		reg = &Registry{}
	})

	Describe("func RegisterSourceDriver()", func() {
		It("registers the source with the given alias", func() {
			expect := sourcedriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterSourceDriver("<alias>", expect)

			d, ok := reg.SourceDriverByAlias("<alias>")
			Expect(ok).To(BeTrue())
			Expect(d).To(Equal(expect))
		})

		It("panics if the alias is already in use", func() {
			d := sourcedriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterSourceDriver("<alias>", d)

			Expect(func() {
				reg.RegisterSourceDriver("<alias>", d)
			}).To(PanicWith("alias is already in use"))
		})
	})

	Describe("func SourceDriverByAlias()", func() {
		It("returns false if there is no source with the given alias", func() {
			_, ok := reg.SourceDriverByAlias("<alias>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func SourceDriverAliases()", func() {
		It("returns a sorted slice of aliases", func() {
			d := sourcedriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterSourceDriver("<b>", d)
			reg.RegisterSourceDriver("<c>", d)
			reg.RegisterSourceDriver("<a>", d)

			Expect(reg.SourceDriverAliases()).To(Equal([]string{
				"<a>",
				"<b>",
				"<c>",
			}))
		})
	})

	When("the registry has a parent", func() {
		var (
			parent           *Registry
			sourceFromParent sourcedriver.Registration
		)

		BeforeEach(func() {
			parent = &Registry{}
			reg.Parent = parent

			sourceFromParent = sourcedriver.Registration{
				Name:        "<name from parent>",
				Description: "<desc from parent>",
			}

			parent.RegisterSourceDriver("<alias>", sourceFromParent)
		})

		Describe("func RegisterSourceDriver()", func() {
			It("shadows a source with the same alias from the parent registry", func() {
				expect := sourcedriver.Registration{
					Name:        "<name>",
					Description: "<desc>",
				}

				reg.RegisterSourceDriver("<alias>", expect)

				d, ok := reg.SourceDriverByAlias("<alias>")
				Expect(ok).To(BeTrue())
				Expect(d).To(Equal(expect))
			})
		})

		Describe("func SourceDriverByAlias()", func() {
			It("falls back to the parent", func() {
				d, ok := reg.SourceDriverByAlias("<alias>")
				Expect(ok).To(BeTrue())
				Expect(d).To(Equal(sourceFromParent))
			})
		})

		Describe("func SourceDriverAliases()", func() {
			It("returns a sorted slice of aliases, including those from the parent", func() {
				d := sourcedriver.Registration{
					Name:        "<name>",
					Description: "<desc>",
				}

				reg.RegisterSourceDriver("<alias>", d) // ensure no dupes with parent

				reg.RegisterSourceDriver("<b>", d)
				reg.RegisterSourceDriver("<c>", d)
				reg.RegisterSourceDriver("<a>", d)

				parent.RegisterSourceDriver("<parent b>", d)
				parent.RegisterSourceDriver("<parent c>", d)
				parent.RegisterSourceDriver("<parent a>", d)

				Expect(reg.SourceDriverAliases()).To(Equal([]string{
					"<a>",
					"<alias>", // note only included once
					"<b>",
					"<c>",
					"<parent a>",
					"<parent b>",
					"<parent c>",
				}))
			})
		})
	})
})
