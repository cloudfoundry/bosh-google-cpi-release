package gsnapshot

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (s GoogleSnapshotService) Find(id string) (*compute.Snapshot, bool, error) {
	s.logger.Debug(googleSnapshotServiceLogTag, "Finding Google Snapshot '%s'", id)
	snapshot, err := s.computeService.Snapshots.Get(s.project, id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.Snapshot{}, false, nil
		}

		return &compute.Snapshot{}, false, bosherr.WrapErrorf(err, "Failed to find Google Snapshot '%s'", id)
	}

	return snapshot, true, nil
}
