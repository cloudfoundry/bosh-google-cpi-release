package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-google-cpi/google/snapshot_service"
)

type DeleteSnapshot struct {
	snapshotService snapshot.Service
}

func NewDeleteSnapshot(
	snapshotService snapshot.Service,
) DeleteSnapshot {
	return DeleteSnapshot{
		snapshotService: snapshotService,
	}
}

func (ds DeleteSnapshot) Run(snapshotCID SnapshotCID) (interface{}, error) {
	err := ds.snapshotService.Delete(string(snapshotCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting snapshot '%s'", snapshotCID)
	}

	return nil, nil
}
