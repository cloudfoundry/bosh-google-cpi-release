package disk

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"bosh-google-cpi/util"
)

// ErrSnapshotPermissionDenied is returned by CreateFromSnapshot when GCP rejects the request
// with a 403, indicating the service account lacks compute.snapshots.useReadOnly permission.
// Callers that catch this error should return Bosh::Clouds::NotSupported so the BOSH director
// falls back to the standard copy-based disk update path.
var ErrSnapshotPermissionDenied = errors.New("compute.snapshots.useReadOnly permission is required to create a disk from snapshot; add the permission to the service account or the disk type change will fall back to the copy path")

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
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusForbidden &&
			strings.Contains(strings.ToLower(gerr.Message), "compute.snapshots.usereadonly") {
			return "", ErrSnapshotPermissionDenied
		}
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk from snapshot")
	}

	if _, err = d.operationService.Waiter(operation, zone, ""); err != nil {
		d.cleanUp(disk.Name)
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk from snapshot")
	}

	return disk.Name, nil
}
