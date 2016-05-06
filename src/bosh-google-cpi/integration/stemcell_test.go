package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {

	It("executes the stemcell lifecycle", func() {
		var stemcellCID string

		By("uploading a stemcell")
		request := fmt.Sprintf(`{
			  "method": "create_stemcell",
			  "arguments": ["", {
				  "name": "bosh-google-kvm-ubuntu-trusty",
				  "version": "3215",
				  "infrastructure": "google",
				  "source_url": "%s"
				}]
			}`, stemcellURL)
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
})
