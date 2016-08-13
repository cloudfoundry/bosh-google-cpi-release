package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Snapshot", func() {
	It("can execute the disk snapshot lifecycle", func() {
		By("creating a disk")
		var diskCID, snapshotCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {"zone": "%v"}, ""]
			}`, zone)
		diskCID = assertSucceedsWithResult(request).(string)

		By("snapshotting the disk")
		request = fmt.Sprintf(`{
			  "method": "snapshot_disk",
			  "arguments": ["%v", {}]
			}`, diskCID)
		snapshotCID = assertSucceedsWithResult(request).(string)

		By("deleting the snapshot")
		request = fmt.Sprintf(`{
			  "method": "delete_snapshot",
			  "arguments": ["%v"]
			}`, snapshotCID)
		assertSucceeds(request)

		By("creating a snapshot with metadata")
		request = fmt.Sprintf(`{
			  "method": "snapshot_disk",
			  "arguments": ["%v", {"deployment": "deployment_name", "job": "job_name", "index": "job_index"}]
			}`, diskCID)
		snapshotCID = assertSucceedsWithResult(request).(string)

		By("deleting the snapshot")
		request = fmt.Sprintf(`{
			  "method": "delete_snapshot",
			  "arguments": ["%v"]
			}`, snapshotCID)
		assertSucceeds(request)

		By("deleting the disk")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)

	})
})
