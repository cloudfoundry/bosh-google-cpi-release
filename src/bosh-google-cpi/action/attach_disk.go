package action

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	disk "bosh-google-cpi/google/disk_service"
	instance "bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/registry"
)

type attachDiskBase struct {
	diskService    disk.Service
	vmService      instance.Service
	registryClient registry.Client
}

type AttachDiskV1 struct {
	attachDiskBase
}

type AttachDiskV2 struct {
	attachDiskBase
}

func NewAttachDiskV1(
	diskService disk.Service,
	vmService instance.Service,
	registryClient registry.Client,
) AttachDiskV1 {
	return AttachDiskV1{
		attachDiskBase{
			diskService:    diskService,
			vmService:      vmService,
			registryClient: registryClient,
		},
	}
}

func NewAttachDiskV2(
	diskService disk.Service,
	vmService instance.Service,
	registryClient registry.Client,
) AttachDiskV2 {
	return AttachDiskV2{
		attachDiskBase{
			diskService:    diskService,
			vmService:      vmService,
			registryClient: registryClient,
		},
	}
}

func (ad1 AttachDiskV1) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	_, err := ad1.attachDiskBase.Run(vmCID, diskCID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ad2 AttachDiskV2) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	diskPath, err := ad2.attachDiskBase.Run(vmCID, diskCID)
	if err != nil {
		return nil, err
	}
	return diskPath, nil
}

func (ad attachDiskBase) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	// Find the disk
	disk, found, err := ad.diskService.Find(string(diskCID), "")
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}
	if !found {
		return nil, api.NewDiskNotFoundError(string(diskCID), false)
	}

	// If this disk isn't already attached to this VM at the IAAS level, do that now
	if !diskAttachedToVM(vmCID, disk.Users) {
		_, err := ad.gceAttach(vmCID, diskCID, disk.SelfLink)
		if err != nil {
			return nil, err
		}
	}

	// The disk may now be configured by the BOSH agent
	dev, err := ad.vmService.DiskDetail(string(vmCID), disk.SelfLink)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Retrieving disk detail before BOSH agent attach")
	}
	err = ad.agentAttach(vmCID, diskCID, dev)

	return dev.Path, err
}

func (ad attachDiskBase) gceAttach(vmCID VMCID, diskCID DiskCID, diskSelfLink string) (*instance.DiskAttachmentDetail, error) {
	// Atach the Disk to the VM
	dev, err := ad.vmService.AttachDisk(string(vmCID), diskSelfLink)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}
	return dev, nil
}

func (ad attachDiskBase) agentAttach(vmCID VMCID, diskCID DiskCID, dev *instance.DiskAttachmentDetail) error {
	// Read VM agent settings
	agentSettings, err := ad.registryClient.Fetch(string(vmCID))
	if err != nil {
		return bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	// Update VM agent settings
	newAgentSettings := agentSettings.AttachPersistentDisk(string(diskCID), dev.Name, dev.Path)
	if err = ad.registryClient.Update(string(vmCID), newAgentSettings); err != nil {
		return bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}
	return nil
}

// This function returns true if the VM has the disk attached.
func diskAttachedToVM(vmCID VMCID, diskUsers []string) bool {
	for _, v := range diskUsers {
		if strings.HasSuffix(v, string(vmCID)) {
			return true
		}
	}
	return false
}
