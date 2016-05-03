package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {

	It("can create a stemcell from a valid JSON request", func() {

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
	})

	It("can delete a stemcell from a valid JSON request", func() {

		request := fmt.Sprintf(`{
			  "method": "delete_stemcell",
			  "arguments": ["%v"]
			}`, stemcellCID)

		response, err := execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.Error).To(BeNil())
		Expect(response.Result).To(BeNil())
	})
})
