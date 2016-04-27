package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Disk", func() {

	It("can create a disk", func() {
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, ""]
			}`)
		diskCID = assertSucceedsWithResult(request).(string)
	})

	It("can attach a disk", func() {
		request := fmt.Sprintf(`{
			  "method": "attach_disk",
			  "arguments": ["%v", "%v"]
			}`, reusableVMName, diskCID)
		assertSucceeds(request)
	})

	It("can confirm the attachment of a disk", func() {
		request := fmt.Sprintf(`{
			  "method": "get_disks",
			  "arguments": ["%v"]
			}`, reusableVMName)
		disks := toStringArray(assertSucceedsWithResult(request).([]interface{}))
		Expect(disks).To(ContainElement(diskCID))
	})

	It("can detach and delete a disk", func() {
		request := fmt.Sprintf(`{
			  "method": "detach_disk",
			  "arguments": ["%v", "%v"]
			}`, reusableVMName, diskCID)
		assertSucceeds(request)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can create and delete a disk in the same zone as an existing machine", func() {
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, "%v"]
			}`, reusableVMName)
		diskCID = assertSucceedsWithResult(request).(string)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can create and delete a PD-SSD disk", func() {
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"type": "pd-ssd"}, "%v"]
			}`, reusableVMName)
		diskCID = assertSucceedsWithResult(request).(string)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})

	It("can create and delete a PD-SSD disk in the same zone as an existing machine", func() {
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"type": "pd-ssd"}, "%v"]
			}`, reusableVMName)
		diskCID = assertSucceedsWithResult(request).(string)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)
	})
})
