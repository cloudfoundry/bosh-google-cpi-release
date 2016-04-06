package snapshot

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (s GoogleSnapshotService) Delete(id string) error {
	snapshot, found, err := s.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "Google Snapshot '%s' does not exists", id)
	}

	if snapshot.Status != googleSnapshotReadyStatus && snapshot.Status != googleSnapshotFailedStatus {
		return bosherr.WrapErrorf(err, "Cannot delete Google Snapshot '%s', status is '%s'", id, snapshot.Status)
	}

	s.logger.Debug(googleSnapshotServiceLogTag, "Deleting Google Snapshot '%s'", id)
	operation, err := s.computeService.Snapshots.Delete(s.project, id).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Snapshot '%s'", id)
	}

	if _, err = s.operationService.Waiter(operation, "", ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Snapshot '%s'", id)
	}

	return nil
}
