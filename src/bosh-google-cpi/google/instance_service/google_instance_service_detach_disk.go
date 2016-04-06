package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
)

func (i GoogleInstanceService) DetachDisk(id string, diskID string) error {
	// Find the instance
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	// Look up for the device name
	var deviceName string
	for _, attachedDisk := range instance.Disks {
		if util.ResourceSplitter(attachedDisk.Source) == diskID {
			deviceName = attachedDisk.DeviceName
		}
	}
	if deviceName == "" {
		return api.NewDiskNotAttachedError(id, diskID, false)
	}

	// Detach the disk
	i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Disk '%s' from Google Instance '%s'", diskID, id)
	operation, err := i.computeService.Instances.DetachDisk(i.project, util.ResourceSplitter(instance.Zone), id, deviceName).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to detach Google Disk '%s' from Google Instance '%s'", diskID, id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to detach Google Disk '%s' from Google Instance '%s'", diskID, id)
	}

	return nil
}
