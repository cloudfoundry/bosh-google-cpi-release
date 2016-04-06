package instance

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

const googleDiskPathPrefix = "/dev/sd"
const googleDiskPathSuffix = "abcdefghijklmnopqrstuvwxyz"

func (i GoogleInstanceService) AttachDisk(id string, diskLink string) (string, string, error) {
	var deviceName, devicePath string

	// Find the instance
	instance, found, err := i.Find(id, "")
	if err != nil {
		return deviceName, devicePath, err
	}
	if !found {
		return deviceName, devicePath, api.NewVMNotFoundError(id)
	}

	deviceName = util.ResourceSplitter(diskLink)
	disk := &compute.AttachedDisk{
		DeviceName: deviceName,
		Mode:       "READ_WRITE",
		Source:     diskLink,
		Type:       "PERSISTENT",
	}

	// Attach the disk
	i.logger.Debug(googleInstanceServiceLogTag, "Attaching Google Disk '%s' to Google Instance '%s'", util.ResourceSplitter(diskLink), id)
	operation, err := i.computeService.Instances.AttachDisk(i.project, util.ResourceSplitter(instance.Zone), id, disk).Do()
	if err != nil {
		return deviceName, devicePath, bosherr.WrapErrorf(err, "Failed to attach Google Disk '%s' to Google Instance '%s'", util.ResourceSplitter(diskLink), id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return deviceName, devicePath, bosherr.WrapErrorf(err, "Failed to attach Google Disk '%s' to Google Instance '%s'", util.ResourceSplitter(diskLink), id)
	}

	// Find the instance again, as we need to get the new attached disks info
	instance, found, err = i.Find(id, "")
	if err != nil {
		return deviceName, devicePath, err
	}
	if !found {
		return deviceName, devicePath, api.NewVMNotFoundError(id)
	}

	// Look up for the device index
	for _, attachedDisk := range instance.Disks {
		if attachedDisk.Source == diskLink {
			deviceIndex := int(attachedDisk.Index)
			devicePath = fmt.Sprintf("%s%s", googleDiskPathPrefix, string(googleDiskPathSuffix[deviceIndex]))
		}
	}

	return deviceName, devicePath, nil
}
