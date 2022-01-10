package source_test

import (
	. "github.com/gritcli/grit/internal/daemon/internal/source"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type List", func() {
	var list List

	BeforeEach(func() {
		list = List{
			{
				Name: "<source>",
			},
		}
	})

	Describe("func ByName()", func() {
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
