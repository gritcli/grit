package config_test

import (
	. "github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/vcsdriver"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type DriverRegistry", func() {
	var reg *DriverRegistry

	BeforeEach(func() {
		reg = &DriverRegistry{}
	})

	Describe("func RegisterVCSDriver()", func() {
		It("registers the VCS with the given alias", func() {
			expect := vcsdriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterVCSDriver("<alias>", expect)

			r, ok := reg.VCSDriverByAlias("<alias>")
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(expect))
		})

		It("panics if the alias is already in use", func() {
			r := vcsdriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterVCSDriver("<alias>", r)

			Expect(func() {
				reg.RegisterVCSDriver("<alias>", r)
			}).To(PanicWith("alias is already in use"))
		})
	})

	Describe("func VCSDriverByAlias()", func() {
		It("returns false if there is no VCS with the given alias", func() {
			_, ok := reg.VCSDriverByAlias("<alias>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func VCSDriverAliases()", func() {
		It("returns a sorted slice of aliases", func() {
			r := vcsdriver.Registration{
				Name:        "<name>",
				Description: "<desc>",
			}

			reg.RegisterVCSDriver("<b>", r)
			reg.RegisterVCSDriver("<c>", r)
			reg.RegisterVCSDriver("<a>", r)

			Expect(reg.VCSDriverAliases()).To(Equal([]string{
				"<a>",
				"<b>",
				"<c>",
			}))
		})
	})

	When("the registry has a parent", func() {
		var (
			parent     *DriverRegistry
			fromParent vcsdriver.Registration
		)

		BeforeEach(func() {
			parent = &DriverRegistry{}
			reg.Parent = parent

			fromParent = vcsdriver.Registration{
				Name:        "<name from parent>",
				Description: "<desc from parent>",
			}

			parent.RegisterVCSDriver("<alias>", fromParent)
		})

		Describe("func RegisterVCSDriver()", func() {
			It("shadows a VCS with the same alias from the parent registry", func() {
				expect := vcsdriver.Registration{
					Name:        "<name>",
					Description: "<desc>",
				}

				reg.RegisterVCSDriver("<alias>", expect)

				r, ok := reg.VCSDriverByAlias("<alias>")
				Expect(ok).To(BeTrue())
				Expect(r).To(Equal(expect))
			})
		})

		Describe("func VCSDriverByAlias()", func() {
			It("falls back to the parent", func() {
				r, ok := reg.VCSDriverByAlias("<alias>")
				Expect(ok).To(BeTrue())
				Expect(r).To(Equal(fromParent))
			})
		})

		Describe("func VCSDrivers()", func() {
			It("returns a map of alias to driver registration, including those from the parent", func() {
				override := vcsdriver.Registration{
					Name:        "<override>",
					Description: "<desc>",
				}

				other := vcsdriver.Registration{
					Name:        "<other>",
					Description: "<desc>",
				}

				reg.RegisterVCSDriver("<alias>", override)
				reg.RegisterVCSDriver("<other>", other)

				Expect(reg.VCSDrivers()).To(Equal(map[string]vcsdriver.Registration{
					"<alias>": override,
					"<other>": other,
				}))
			})
		})

		Describe("func VCSDriverAliases()", func() {
			It("returns a sorted slice of aliases, including those from the parent", func() {
				r := vcsdriver.Registration{
					Name:        "<name>",
					Description: "<desc>",
				}

				reg.RegisterVCSDriver("<alias>", r) // ensure no dupes with parent

				reg.RegisterVCSDriver("<b>", r)
				reg.RegisterVCSDriver("<c>", r)
				reg.RegisterVCSDriver("<a>", r)

				parent.RegisterVCSDriver("<parent b>", r)
				parent.RegisterVCSDriver("<parent c>", r)
				parent.RegisterVCSDriver("<parent a>", r)

				Expect(reg.VCSDriverAliases()).To(Equal([]string{
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
