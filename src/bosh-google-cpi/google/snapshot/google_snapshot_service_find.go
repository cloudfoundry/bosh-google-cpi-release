package snapshot

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"google.golang.org/api/googleapi"
)

func (s GoogleSnapshotService) Find(id string) (Snapshot, bool, error) {
	s.logger.Debug(googleSnapshotServiceLogTag, "Finding Google Snapshot '%s'", id)
	snapshotItem, err := s.computeService.Snapshots.Get(s.project, id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return Snapshot{}, false, nil
		}

		return Snapshot{}, false, bosherr.WrapErrorf(err, "Failed to find Google Snapshot '%s'", id)
	}

	snapshot := Snapshot{
		Name:     snapshotItem.Name,
		SelfLink: snapshotItem.SelfLink,
		Status:   snapshotItem.Status,
	}
	return snapshot, true, nil
}
