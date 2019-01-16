package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"
)

var _ = Describe("Info", func() {
	var subject Info

	BeforeEach(func() {
		subject = NewInfo()
	})

	Describe("Run", func() {
		var response InfoResult

		BeforeEach(func() {
			var err error

			response, err = subject.Run()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("stemcell_formats", func() {
			It("supports google-light", func() {
				Expect(response.StemcellFormats).To(ContainElement("google-light"))
			})

			It("supports google-rawdisk", func() {
				Expect(response.StemcellFormats).To(ContainElement("google-rawdisk"))
			})
		})

		Context("api_version", func() {
			It("returns the latest api_version", func() {
				Expect(response.ApiVersion).To(Equal(1))
			})
		})

	})
})
