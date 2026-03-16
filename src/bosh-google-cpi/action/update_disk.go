package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk"
	"bosh-google-cpi/google/disktype"
	"bosh-google-cpi/google/snapshot"
	"bosh-google-cpi/util"
)

type UpdateDisk struct {
	diskService     disk.Service
	diskTypeService disktype.Service
	snapshotService snapshot.Service
}

func NewUpdateDisk(
	diskService disk.Service,
	diskTypeService disktype.Service,
	snapshotService snapshot.Service,
) UpdateDisk {
	return UpdateDisk{
		diskService:     diskService,
		diskTypeService: diskTypeService,
		snapshotService: snapshotService,
	}
}

func (ud UpdateDisk) Run(diskCID DiskCID, newSize int, cloudProps DiskCloudProperties) (interface{}, error) {
	// Find the existing disk
	existingDisk, found, err := ud.diskService.Find(string(diskCID), "")
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Updating disk '%s'", diskCID)
	}
	if !found {
		return nil, api.NewDiskNotFoundError(string(diskCID), false)
	}

	zone := existingDisk.Zone
	sizeGib := util.ConvertMib2Gib(newSize)
	if sizeGib < int(existingDisk.SizeGb) {
		sizeGib = int(existingDisk.SizeGb)
	}

	// Resolve the target disk type; default to the existing type so we never
	// silently change it when no type is specified in cloud_properties.
	diskTypeSelfLink := existingDisk.Type
	if cloudProps.DiskType != "" {
		dt, found, err := ud.diskTypeService.Find(cloudProps.DiskType, zone)
		if err != nil {
			return nil, bosherr.WrapErrorf(err, "Updating disk '%s': finding disk type '%s'", diskCID, cloudProps.DiskType)
		}
		if !found {
			return nil, bosherr.Errorf("Updating disk '%s': disk type '%s' not found", diskCID, cloudProps.DiskType)
		}
		diskTypeSelfLink = dt.SelfLink
	}

	if existingDisk.Type == diskTypeSelfLink {
		if sizeGib <= int(existingDisk.SizeGb) {
			return string(diskCID), nil // already at target size and type
		}
		if err := ud.diskService.Resize(string(diskCID), sizeGib); err != nil {
			return nil, bosherr.WrapErrorf(err, "Updating disk '%s': resizing", diskCID)
		}
		return string(diskCID), nil
	}

	// Snapshot the existing disk
	snapshotID, err := ud.snapshotService.Create(string(diskCID), "update-disk snapshot", zone)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Updating disk '%s': creating snapshot", diskCID)
	}

	// Find the snapshot to get its SelfLink
	snap, found, err := ud.snapshotService.Find(snapshotID)
	if err != nil {
		_ = ud.snapshotService.Delete(snapshotID) //nolint:errcheck
		return nil, bosherr.WrapErrorf(err, "Updating disk '%s': finding snapshot '%s'", diskCID, snapshotID)
	}
	if !found {
		_ = ud.snapshotService.Delete(snapshotID) //nolint:errcheck
		return nil, bosherr.Errorf("Updating disk '%s': snapshot '%s' not found after creation", diskCID, snapshotID)
	}

	// Create the new disk from snapshot BEFORE deleting the old one.
	// This ensures data is never lost: if creation fails, the old disk
	// is still intact and the caller can retry or fall back.
	newDiskID, err := ud.diskService.CreateFromSnapshot(snap.SelfLink, sizeGib, diskTypeSelfLink, zone)
	if err != nil {
		_ = ud.snapshotService.Delete(snapshotID) //nolint:errcheck
		return nil, bosherr.WrapErrorf(err, "Updating disk '%s': recreating from snapshot '%s'", diskCID, snapshotID)
	}

	if err := ud.diskService.Delete(string(diskCID)); err != nil {
		_ = ud.snapshotService.Delete(snapshotID) //nolint:errcheck
		return newDiskID, nil
	}

	_ = ud.snapshotService.Delete(snapshotID) //nolint:errcheck

	return newDiskID, nil
}
