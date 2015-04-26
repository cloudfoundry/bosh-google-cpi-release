package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/snapshot"
)

type DeleteSnapshot struct {
	snapshotService gsnapshot.GoogleSnapshotService
}

func NewDeleteSnapshot(
	snapshotService gsnapshot.GoogleSnapshotService,
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
