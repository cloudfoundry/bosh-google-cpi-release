package disk

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (d GoogleDiskService) Resize(id string, newSize int) error {
	newsizeGB := &compute.DisksResizeRequest{
		SizeGb: int64(newSize),
	}

	disk, found, err := d.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewDiskNotFoundError(id, false)
	}

	if disk.Status != googleDiskReadyStatus && disk.Status != googleDiskFailedStatus {
		return bosherr.WrapErrorf(err, "Cannot resize Google Disk '%s', status is '%s'", id, disk.Status)
	}

	d.logger.Debug(googleDiskServiceLogTag, "Resizing Google Disk '%s'", id)
	operation, err := d.computeService.Disks.Resize(d.project, util.ResourceSplitter(disk.Zone), id, newsizeGB).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to resize Google Disk '%s'", id)
	}

	if _, err = d.operationService.Waiter(operation, disk.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to resize Google Disk '%s'", id)
	}

	return nil
}

//TODO:
// cannot resize to smaller disk (this should be handled already by the director call)
// skip resize when values are the same
// cannot resize as disk is still attached.
