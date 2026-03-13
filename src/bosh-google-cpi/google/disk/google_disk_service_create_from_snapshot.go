package disk

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"google.golang.org/api/compute/v1"

	"bosh-google-cpi/util"
)

// CreateFromSnapshot creates a new disk in zone from snapshotSelfLink with the given size (GiB).
// diskType must be the full SelfLink URL of the desired disk type (e.g. the value returned by
// DiskTypeService.Find). If empty, GCP uses the zone's default disk type.
func (d GoogleDiskService) CreateFromSnapshot(snapshotSelfLink string, size int, diskType string, zone string) (string, error) {
	uuidStr, err := d.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random Google Disk name")
	}

	disk := &compute.Disk{
		Name:           fmt.Sprintf("%s-%s", googleDiskNamePrefix, uuidStr),
		Description:    googleDiskDescription,
		SizeGb:         int64(size),
		SourceSnapshot: snapshotSelfLink,
	}

	if diskType != "" {
		disk.Type = diskType
	}

	d.logger.Debug(googleDiskServiceLogTag, "Creating Google Disk from snapshot with params: %#v", disk)
	operation, err := d.computeService.Disks.Insert(d.project, util.ResourceSplitter(zone), disk).Do()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk from snapshot")
	}

	if _, err = d.operationService.Waiter(operation, zone, ""); err != nil {
		d.cleanUp(disk.Name)
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk from snapshot")
	}

	return disk.Name, nil
}
