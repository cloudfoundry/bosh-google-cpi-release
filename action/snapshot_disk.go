package action

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/disk"
	"github.com/frodenas/bosh-google-cpi/google/snapshot"
)

type CreateSnapshot struct {
	snapshotService gsnapshot.GoogleSnapshotService
	diskService     gdisk.GoogleDiskService
}

func NewSnapshotDisk(
	snapshotService gsnapshot.GoogleSnapshotService,
	diskService gdisk.GoogleDiskService,
) CreateSnapshot {
	return CreateSnapshot{
		snapshotService: snapshotService,
		diskService:     diskService,
	}
}

func (sd CreateSnapshot) Run(diskCID DiskCID, metadata SnapshotMetadata) (SnapshotCID, error) {
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
