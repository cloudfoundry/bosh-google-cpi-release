package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Snapshot", func() {

	It("can create a snapshot", func() {
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, ""]
			}`)
		diskCID = assertSucceedsWithResult(request).(string)
		request = fmt.Sprintf(`{
			  "method": "snapshot_disk",
			  "arguments": ["%v", {}]
			}`, diskCID)
		snapshotCID = assertSucceedsWithResult(request).(string)
	})

	It("can delete a snapshot", func() {
		request := fmt.Sprintf(`{
			  "method": "delete_snapshot",
			  "arguments": ["%v"]
			}`, snapshotCID)
		assertSucceeds(request)
	})

	It("can create and delete a snapshot with metadata", func() {
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, ""]
			}`)
		diskCID = assertSucceedsWithResult(request).(string)
		request = fmt.Sprintf(`{
			  "method": "snapshot_disk",
			  "arguments": ["%v", {"deployment": "deployment_name", "job": "job_name", "index": "job_index"}]
			}`, diskCID)
		snapshotCID = assertSucceedsWithResult(request).(string)

	})

})
