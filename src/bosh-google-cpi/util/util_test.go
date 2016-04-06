package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/util"
)

var _ = Describe("Util", func() {
	Describe("ConvertMib2Gib", func() {
		It("converts Mib to Gib", func() {
			Expect(ConvertMib2Gib(32768)).To(Equal(32))
		})
	})

	Describe("ResourceSplitter", func() {
		It("splits the resource name", func() {
			Expect(ResourceSplitter("prefix/fake-resource")).To(Equal("fake-resource"))
		})

		Context("when resource does not contains slash", func() {
			It("returns the string", func() {
				Expect(ResourceSplitter("fake-resource")).To(Equal("fake-resource"))
			})
		})
	})
})
