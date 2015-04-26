package gdisk

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
)

func (d GoogleDiskService) Delete(id string) error {
	disk, found, err := d.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Disk '%s' not found", id)
	}

	if disk.Status != googleDiskReadyStatus {
		return bosherr.WrapErrorf(err, "Cannot delete Google Disk '%s', status is '%s'", id, disk.Status)
	}

	d.logger.Debug(googleDiskServiceLogTag, "Deleting Google Disk '%s'", id)
	operation, err := d.computeService.Disks.Delete(d.project, gutil.ResourceSplitter(disk.Zone), id).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Disk '%s'", id)
	}

	if _, err = d.operationService.Waiter(operation, disk.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Disk '%s'", id)
	}

	return nil
}
