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

func (i GoogleInstanceService) AttachDisk(id string, diskLink string) (*DiskAttachmentDetail, error) {
	// Find the instance
	instance, found, err := i.Find(id, "")
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, api.NewVMNotFoundError(id)
	}

	deviceName := i.diskDeviceName(diskLink)
	disk := &compute.AttachedDisk{
		DeviceName: deviceName,
		Mode:       "READ_WRITE",
		Source:     diskLink,
		Type:       "PERSISTENT",
	}

	// Attach the disk
	i.logger.Debug(googleInstanceServiceLogTag, "Attaching Google Disk '%s' to Google Instance '%s'", deviceName, id)
	operation, err := i.computeService.Instances.AttachDisk(i.project, util.ResourceSplitter(instance.Zone), id, disk).Do()
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Failed to attach Google Disk '%s' to Google Instance '%s'", deviceName, id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return nil, bosherr.WrapErrorf(err, "Failed to attach Google Disk '%s' to Google Instance '%s'", deviceName, id)
	}

	return i.DiskDetail(id, diskLink)
}

func (i GoogleInstanceService) DiskDetail(vmID string, diskLink string) (*DiskAttachmentDetail, error) {
	deviceName := i.diskDeviceName(diskLink)

	// Find the instance the disk should be attached to
	instance, found, err := i.Find(vmID, "")
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, api.NewVMNotFoundError(vmID)
	}

	// Derive the disk's attachment path based on its index
	for _, attachedDisk := range instance.Disks {
		if attachedDisk.Source == diskLink {
			deviceIndex := int(attachedDisk.Index)
			dev := &DiskAttachmentDetail{
				Name: deviceName,
				Path: fmt.Sprintf("%s%s", googleDiskPathPrefix, string(googleDiskPathSuffix[deviceIndex])),
			}
			return dev, nil
		}
	}

	return nil, bosherr.Errorf("Disk %q is not attached to instance %q", diskLink, vmID)
}

func (i GoogleInstanceService) diskDeviceName(diskLink string) string {
	return util.ResourceSplitter(diskLink)
}
