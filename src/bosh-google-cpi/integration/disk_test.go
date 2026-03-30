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

		By("attaching the disk again without failing")
		request = fmt.Sprintf(`{
			  "method": "attach_disk",
			  "arguments": ["%v", "%v"]
			}`, vmCID, diskCID)
		assertSucceeds(request)

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

	It("can create and resize and delete a disk", func() {
		By("creating a disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("resizing the disk")
		request = fmt.Sprintf(`{
			  "method": "resize_disk",
			  "arguments": ["%v", 52768]
			}`, diskCID)
		assertSucceeds(request)

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can update a disk with the same type and size (no-op)", func() {
		By("creating a pd-ssd disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"type": "pd-ssd", "zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("updating with the same type and size")
		request = fmt.Sprintf(`{
			  "method": "update_disk",
			  "arguments": ["%v", 32768, {"type": "pd-ssd"}]
			}`, diskCID)
		result := assertSucceedsWithResult(request).(string)
		Expect(result).To(Equal(diskCID))

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can update a disk by resizing in-place when only size changes", func() {
		By("creating a pd-ssd disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"type": "pd-ssd", "zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("updating with same type but larger size")
		request = fmt.Sprintf(`{
			  "method": "update_disk",
			  "arguments": ["%v", 52768, {"type": "pd-ssd"}]
			}`, diskCID)
		result := assertSucceedsWithResult(request).(string)
		Expect(result).To(Equal(diskCID))

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("fails to update a disk to a smaller size", func() {
		By("creating a disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"type": "pd-ssd", "zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("attempting to shrink the disk")
		request = fmt.Sprintf(`{
			  "method": "update_disk",
			  "arguments": ["%v", 10240, {"type": "pd-ssd"}]
			}`, diskCID)
		err := assertFails(request)
		Expect(err.Error()).To(ContainSubstring("smaller than current size"))

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can update a disk by changing its type via snapshot", func() {
		By("creating a pd-ssd disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"type": "pd-ssd", "zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("updating to hyperdisk-balanced type")
		request = fmt.Sprintf(`{
			  "method": "update_disk",
			  "arguments": ["%v", 32768, {"type": "hyperdisk-balanced"}]
			}`, diskCID)
		newDiskCID := assertSucceedsWithResult(request).(string)
		Expect(newDiskCID).ToNot(Equal(diskCID))

		By("confirming the new disk exists")
		request = fmt.Sprintf(`{
			  "method": "has_disk",
			  "arguments": ["%v"]
			}`, newDiskCID)
		found := assertSucceedsWithResult(request).(bool)
		Expect(found).To(BeTrue())

		By("confirming the old disk was deleted")
		request = fmt.Sprintf(`{
			  "method": "has_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		found = assertSucceedsWithResult(request).(bool)
		Expect(found).To(BeFalse())

		By("deleting the new disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, newDiskCID)
		assertSucceeds(request)
	})

})
