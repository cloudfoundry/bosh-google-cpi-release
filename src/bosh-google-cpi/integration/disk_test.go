package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Disk", func() {

	It("executes the disk lifecycle", func() {
		By("creating a disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("confirming a disk exists")
		request = fmt.Sprintf(`{
			  "method": "has_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)

		By("creating a VM")
		var vmCID string
		request = fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "zone": "%v"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, zone, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("attaching the disk")
		request = fmt.Sprintf(`{
			  "method": "attach_disk",
			  "arguments": ["%v", "%v"]
			}`, vmCID, diskCID)
		assertSucceeds(request)

		By("confirming the attachment of a disk")
		request = fmt.Sprintf(`{
			  "method": "get_disks",
			  "arguments": ["%v"]
			}`, vmCID)
		disks := toStringArray(assertSucceedsWithResult(request).([]interface{}))
		Expect(disks).To(ContainElement(diskCID))

		By("detaching and deleting a disk")
		request = fmt.Sprintf(`{
			  "method": "detach_disk",
			  "arguments": ["%v", "%v"]
			}`, vmCID, diskCID)
		assertSucceeds(request)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)

		By("confirming a disk does not exist")
		request = fmt.Sprintf(`{
			  "method": "has_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		found := assertSucceedsWithResult(request).(bool)
		Expect(found).To(BeFalse())

		By("creating a disk in the same zone as a VM")
		request = fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, "%v"]
			}`, vmCID)
		diskCID = assertSucceedsWithResult(request).(string)

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create and delete a PD-SSD disk", func() {
		By("creating a disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"type": "pd-ssd", "zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can create and delete a 1TB PD-SSD disk", func() {
		By("creating a disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [1024000, {"type": "pd-ssd", "zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})
})
