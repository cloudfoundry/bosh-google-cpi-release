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

	// If the disk already has the target type and meets the size requirement, skip.
	if existingDisk.Type == diskTypeSelfLink && sizeGib <= int(existingDisk.SizeGb) {
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
		_ = ud.snapshotService.Delete(snapshotID)
		return nil, bosherr.WrapErrorf(err, "Updating disk '%s': finding snapshot '%s'", diskCID, snapshotID)
	}
	if !found {
		_ = ud.snapshotService.Delete(snapshotID)
		return nil, bosherr.Errorf("Updating disk '%s': snapshot '%s' not found after creation", diskCID, snapshotID)
	}

	// Delete the old disk
	if err := ud.diskService.Delete(string(diskCID)); err != nil {
		// Snapshot exists, but we failed to delete the old disk. Clean up snapshot.
		_ = ud.snapshotService.Delete(snapshotID)
		return nil, bosherr.WrapErrorf(err, "Updating disk '%s': deleting old disk", diskCID)
	}

	newDiskID, err := ud.diskService.CreateFromSnapshot(snap.SelfLink, sizeGib, diskTypeSelfLink, zone)
	if err != nil {
		// Old disk is gone; snapshot remains for manual recovery.
		return nil, bosherr.WrapErrorf(err, "Updating disk '%s': recreating from snapshot '%s'", diskCID, snapshotID)
	}

	// Non-fatal: disk was recreated successfully, snapshot is just orphaned.
	_ = ud.snapshotService.Delete(snapshotID)

	return newDiskID, nil
}
