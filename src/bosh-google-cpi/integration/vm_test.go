package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM", func() {

	BeforeEach(func() {
		Expect(googleProject).ToNot(Equal(""), "GOOGLE_PROJECT must be set")
		Expect(externalStaticIP).ToNot(Equal(""), "EXTERNAL_STATIC_IP must be set")
	})

	It("can create a simple VM", func() {
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
	})

	It("can confirm existence of an existing VM", func() {
		request := fmt.Sprintf(`{
			  "method": "has_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		exists := assertSucceedsWithResult(request).(bool)
		Expect(exists).To(Equal(true))
	})

	It("can reboot an existing VM", func() {
		request := fmt.Sprintf(`{
			  "method": "reboot_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a disk", func() {
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, ""]
			}`)
		diskCID = assertSucceedsWithResult(request).(string)
	})

	It("can attach a disk to a VM", func() {
		request := fmt.Sprintf(`{
			  "method": "attach_disk",
			  "arguments": ["%v", "%v"]
			}`, vmCID, diskCID)
		assertSucceeds(request)
	})

	It("can confirm the attachment of a disk to a VM", func() {
		request := fmt.Sprintf(`{
			  "method": "get_disks",
			  "arguments": ["%v"]
			}`, vmCID)
		disks := toStringArray(assertSucceedsWithResult(request).([]interface{}))
		Expect(disks).To(ContainElement(diskCID))
	})

	It("can detach and delete a disk from a VM", func() {
		request := fmt.Sprintf(`{
			  "method": "detach_disk",
			  "arguments": ["%v", "%v"]
			}`, vmCID, diskCID)
		assertSucceeds(request)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can delete a VM", func() {
		request := fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create and delete a custom VM", func() {
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "cpu": 2,
				  "ram": 5120
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})
})
