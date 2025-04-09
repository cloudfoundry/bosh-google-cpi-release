package disk

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
)

func (d GoogleDiskService) Delete(id string) error {
	disk, found, err := d.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewDiskNotFoundError(id, false)
	}

	if disk.Status != googleDiskReadyStatus && disk.Status != googleDiskFailedStatus {
		return bosherr.WrapErrorf(err, "Cannot delete Google Disk '%s', status is '%s'", id, disk.Status)
	}

	d.logger.Debug(googleDiskServiceLogTag, "Deleting Google Disk '%s'", id)
	operation, err := d.computeService.Disks.Delete(d.project, util.ResourceSplitter(disk.Zone), id).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Disk '%s'", id)
	}

	if _, err = d.operationService.Waiter(operation, disk.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Disk '%s'", id)
	}

	return nil
}
