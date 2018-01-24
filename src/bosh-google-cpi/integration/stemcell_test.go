package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {
	It("executes the stemcell lifecycle with an image_url", func() {
		var stemcellCID string

		By("uploading a stemcell with image_url")
		request := fmt.Sprintf(`{
         "method": "create_stemcell",
         "arguments": ["", {
           "name": "bosh-google-kvm-ubuntu-trusty",
           "version": "3215",
           "infrastructure": "google",
           "image_url": "%s"
         }]
       }`, imageURL)
		stemcellCID = assertSucceedsWithResult(request).(string)

		By("deleting the stemcell")
		request = fmt.Sprintf(`{
         "method": "delete_stemcell",
         "arguments": ["%v"]
       }`, stemcellCID)

		response, err := execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.Error).To(BeNil())
		Expect(response.Result).To(BeNil())
	})

	It("quietly ignores missing stemcells", func() {
		By("deleting the stemcell")
		request := `{
         "method": "delete_stemcell",
         "arguments": ["stemcell-ab5f5573-62d5-4f72-6620-3ab8f1aed3cf"]
		}`

		response, err := execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.Error).To(BeNil())
		Expect(response.Result).To(BeNil())
	})
	//TODO: Add heavy stemcell test
})
