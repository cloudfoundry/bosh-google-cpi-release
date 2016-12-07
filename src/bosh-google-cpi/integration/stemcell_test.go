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
				  "source_url": "%s",
				  "raw_disk_sha1": "%s"
				}]
			}`, stemcellURL, stemcellSHA1)
		stemcellCID = assertSucceedsWithResult(request).(string)

		By("deleting the stemcell")
		request = fmt.Sprintf(`{
			  "method": "delete_stemcell",
			  "arguments": ["%v"]
			}`, stemcellCID)
		assertSucceeds(request)
	})

	It("executes the stemcell lifecycle without a checksum for backwards compatibility", func() {
		var stemcellCID string

		By("uploading a stemcell without a checksum")
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
		assertSucceeds(request)
	})

	//This currently fails due to issue: https://github.com/google/google-api-go-client/issues/181
	XIt("returns an error if the checksum does not match", func() {
		By("uploading a stemcell with an invalid checksum")
		request := fmt.Sprintf(`{
			  "method": "create_stemcell",
			  "arguments": ["", {
				  "name": "bosh-google-kvm-ubuntu-trusty",
				  "version": "3215",
				  "infrastructure": "google",
				  "source_url": "%s",
				  "raw_disk_sha1": "my-invalid-checksum"
				}]
			}`, stemcellURL)
		err := assertFails(request)
		Expect(err.Error()).To(ContainSubstring("TODO"))
	})

	//TODO: Add heavy stemcell test
})
