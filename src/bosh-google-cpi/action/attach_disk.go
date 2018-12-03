package action

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/instance_service"

	"bosh-google-cpi/registry"
)

type AttachDisk struct {
	diskService    disk.Service
	vmService      instance.Service
	registryClient registry.Client
}

func NewAttachDisk(
	diskService disk.Service,
	vmService instance.Service,
	registryClient registry.Client,
) AttachDisk {
	return AttachDisk{
		diskService:    diskService,
		vmService:      vmService,
		registryClient: registryClient,
	}
}

func (ad AttachDisk) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	// Find the disk
	disk, found, err := ad.diskService.Find(string(diskCID), "")
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}
	if !found {
		return nil, api.NewDiskNotFoundError(string(diskCID), false)
	}

	// If the persistent disk is already attached to the VM, only attempt the
	// BOSH agent attachment.
	if diskAttachedToVM(vmCID, disk.Users) {
		return nil, nil
	}

	// If this disk is attached to another VM, attaching here will fail so return
	// a useful error before trying.
	if len(disk.Users) > 0 {
		return nil, api.NewDiskAlreadyAttachedError(string(diskCID), disk.Users, false)
	}

	err, deviceName, devicePath := ad.gceAttach(vmCID, diskCID, disk.SelfLink)
	if err != nil {
		return nil, err
	}

	err = ad.agentAttach(vmCID, diskCID, deviceName, devicePath)
	return nil, err
}

func (ad AttachDisk) gceAttach(vmCID VMCID, diskCID DiskCID, diskSelfLink string) (error, string, string) {
	// Atach the Disk to the VM
	deviceName, devicePath, err := ad.vmService.AttachDisk(string(vmCID), diskSelfLink)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return err, deviceName, devicePath
		}
		return bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID), deviceName, devicePath
	}
	return nil, deviceName, devicePath
}

func (ad AttachDisk) agentAttach(vmCID VMCID, diskCID DiskCID, deviceName string, devicePath string) error {
	// Read VM agent settings
	agentSettings, err := ad.registryClient.Fetch(string(vmCID))
	if err != nil {
		return bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	// Update VM agent settings
	newAgentSettings := agentSettings.AttachPersistentDisk(string(diskCID), deviceName, devicePath)
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
