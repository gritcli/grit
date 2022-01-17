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

			r, ok := reg.SourceDriverByAlias("<alias>")
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(expect))
		})

		It("panics if the alias is already in use", func() {
			r := sourcedriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterSourceDriver("<alias>", r)

			Expect(func() {
				reg.RegisterSourceDriver("<alias>", r)
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
			r := sourcedriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterSourceDriver("<b>", r)
			reg.RegisterSourceDriver("<c>", r)
			reg.RegisterSourceDriver("<a>", r)

			Expect(reg.SourceDriverAliases()).To(Equal([]string{
				"<a>",
				"<b>",
				"<c>",
			}))
		})
	})

	When("the registry has a parent", func() {
		var (
			parent     *Registry
			fromParent sourcedriver.Registration
		)

		BeforeEach(func() {
			parent = &Registry{}
			reg.Parent = parent

			fromParent = sourcedriver.Registration{
				Name:        "<name from parent>",
				Description: "<desc from parent>",
			}

			parent.RegisterSourceDriver("<alias>", fromParent)
		})

		Describe("func RegisterSourceDriver()", func() {
			It("shadows a source with the same alias from the parent registry", func() {
				expect := sourcedriver.Registration{
					Name:        "<name>",
					Description: "<desc>",
				}

				reg.RegisterSourceDriver("<alias>", expect)

				r, ok := reg.SourceDriverByAlias("<alias>")
				Expect(ok).To(BeTrue())
				Expect(r).To(Equal(expect))
			})
		})

		Describe("func SourceDriverByAlias()", func() {
			It("falls back to the parent", func() {
				r, ok := reg.SourceDriverByAlias("<alias>")
				Expect(ok).To(BeTrue())
				Expect(r).To(Equal(fromParent))
			})
		})

		Describe("func SourceDrivers()", func() {
			It("returns a map of alias to driver registration, including those from the parent", func() {
				override := sourcedriver.Registration{
					Name:        "<override>",
					Description: "<desc>",
				}

				other := sourcedriver.Registration{
					Name:        "<other>",
					Description: "<desc>",
				}

				reg.RegisterSourceDriver("<alias>", override)
				reg.RegisterSourceDriver("<other>", other)

				Expect(reg.SourceDrivers()).To(Equal(map[string]sourcedriver.Registration{
					"<alias>": override,
					"<other>": other,
				}))
			})
		})

		Describe("func SourceDriverAliases()", func() {
			It("returns a sorted slice of aliases, including those from the parent", func() {
				r := sourcedriver.Registration{
					Name:        "<name>",
					Description: "<desc>",
				}

				reg.RegisterSourceDriver("<alias>", r) // ensure no dupes with parent

				reg.RegisterSourceDriver("<b>", r)
				reg.RegisterSourceDriver("<c>", r)
				reg.RegisterSourceDriver("<a>", r)

				parent.RegisterSourceDriver("<parent b>", r)
				parent.RegisterSourceDriver("<parent c>", r)
				parent.RegisterSourceDriver("<parent a>", r)

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
