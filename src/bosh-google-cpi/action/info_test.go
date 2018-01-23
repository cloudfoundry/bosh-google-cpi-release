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
		Context("stemcell_formats", func() {
			var response InfoResult

			BeforeEach(func() {
				var err error

				response, err = subject.Run()
				Expect(err).NotTo(HaveOccurred())
			})

			It("supports google-light", func() {
				Expect(response.StemcellFormats).To(ContainElement("google-light"))
			})

			It("supports google-rawdisk", func() {
				Expect(response.StemcellFormats).To(ContainElement("google-rawdisk"))
			})
		})
	})
})
