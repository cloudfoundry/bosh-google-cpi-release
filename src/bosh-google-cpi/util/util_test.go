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

	Describe("RegionFromZone", func() {
		It("successfully parses region from well-formed zone", func() {
			Expect(RegionFromZone("us-west2-c")).To(Equal("us-west2"))
		})

		It("fails to parse region from mal-formed zone", func() {
			Expect(RegionFromZone("notAZone")).To(Equal(""))
		})
	})

	Describe("ZoneFromURL", func() {
		It("successfully parses zone from well-formed URL", func() {
			Expect(ZoneFromURL("https://www.googleapis.com/compute/v1/projects/test-project-id/zones/us-west2-a")).To(Equal("us-west2-a"))
		})

		It("failed to parse zone from mal-formed URL", func() {
			Expect(ZoneFromURL("https://www.googleapis.com/compute/v1/projects/test-project-id/")).To(Equal(""))
		})
	})

	Describe("RegionFromURL", func() {
		It("successfully parses region from well-formed URL", func() {
			Expect(RegionFromURL("https://www.googleapis.com/compute/v1/projects/test-project-id/regions/us-central1")).To(Equal("us-central1"))
		})

		It("fails to parse region from mal-formed URL", func() {
			Expect(RegionFromURL("https://www.googleapis.com/compute/v1/projects/test-project-id/")).To(Equal(""))
		})
	})
})
