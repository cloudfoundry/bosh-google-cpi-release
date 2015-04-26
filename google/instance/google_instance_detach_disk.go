package ginstance

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
)

func (i GoogleInstanceService) DetachDisk(id string, diskId string) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Instance '%s' not found", id)
	}

	var deviceName string
	for _, attachedDisk := range instance.Disks {
		if gutil.ResourceSplitter(attachedDisk.Source) == diskId {
			deviceName = attachedDisk.DeviceName
		}
	}
	if deviceName == "" {
		return bosherr.Errorf("Google Disk '%s' is not attached to Google Instance '%s'", diskId, id)
	}

	i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Disk '%s' from Google Instance '%s'", diskId, id)
	operation, err := i.computeService.Instances.DetachDisk(i.project, gutil.ResourceSplitter(instance.Zone), id, deviceName).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to detach Google Disk '%s' from Google Instance '%s'", diskId, id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to detach Google Disk '%s' from Google Instance '%s'", diskId, id)
	}

	return nil
}
