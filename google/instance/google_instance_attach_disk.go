package ginstance

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceService) AttachDisk(id string, diskLink string) (string, error) {
	// Find the instance
	instance, found, err := i.Find(id, "")
	if err != nil {
		return "", err
	}
	if !found {
		return "", api.NewVMNotFoundError(id)
	}

	deviceName := gutil.ResourceSplitter(diskLink)
	disk := &compute.AttachedDisk{
		DeviceName: deviceName,
		Mode:       "READ_WRITE",
		Source:     diskLink,
		Type:       "PERSISTENT",
	}

	// Attach the disk
	i.logger.Debug(googleInstanceServiceLogTag, "Attaching Google Disk '%s' to Google Instance '%s'", gutil.ResourceSplitter(diskLink), id)
	operation, err := i.computeService.Instances.AttachDisk(i.project, gutil.ResourceSplitter(instance.Zone), id, disk).Do()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to attach Google Disk '%s' to Google Instance '%s'", gutil.ResourceSplitter(diskLink), id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to attach Google Disk '%s' to Google Instance '%s'", gutil.ResourceSplitter(diskLink), id)
	}

	return deviceName, nil
}
