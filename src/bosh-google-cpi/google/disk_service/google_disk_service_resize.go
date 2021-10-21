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

	if disk.SizeGb == newsizeGB.SizeGb {
		d.logger.Debug(googleDiskServiceLogTag, "Skipping resize Google Disk '%s', becasue current value '%#v'is equal to new value '%#v'", id, disk.SizeGb, newsizeGB)
	} else if disk.SizeGb > newsizeGB.SizeGb {
		return bosherr.WrapErrorf(err, "Skipping resize Google Disk '%s', cannot resize volume to a smaller size from '%#v' to '%#v'", id, disk.SizeGb, newsizeGB)
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
