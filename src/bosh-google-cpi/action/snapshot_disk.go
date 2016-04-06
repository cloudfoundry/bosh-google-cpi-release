package action

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/snapshot_service"
)

type SnapshotDisk struct {
	snapshotService snapshot.Service
	diskService     disk.Service
}

func NewSnapshotDisk(
	snapshotService snapshot.Service,
	diskService disk.Service,
) SnapshotDisk {
	return SnapshotDisk{
		snapshotService: snapshotService,
		diskService:     diskService,
	}
}

func (sd SnapshotDisk) Run(diskCID DiskCID, metadata SnapshotMetadata) (SnapshotCID, error) {
	var description string

	// Find the disk
	disk, found, err := sd.diskService.Find(string(diskCID), "")
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to find disk '%s'", diskCID)
	}
	if !found {
		return "", api.NewDiskNotFoundError(string(diskCID), false)
	}

	// Create the disk snapshot
	if metadata.Deployment != "" && metadata.Job != "" && metadata.Index != "" {
		description = fmt.Sprintf("%s/%s/%s", metadata.Deployment, metadata.Job, metadata.Index)
	}

	snapshot, err := sd.snapshotService.Create(string(diskCID), description, disk.Zone)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating disk snapshot")
	}

	return SnapshotCID(snapshot), nil
}
