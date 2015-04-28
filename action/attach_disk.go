package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/disk"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/registry"
)

type AttachDisk struct {
	diskService    gdisk.GoogleDiskService
	vmService      ginstance.GoogleInstanceService
	registryClient registry.Client
}

func NewAttachDisk(
	diskService gdisk.GoogleDiskService,
	vmService ginstance.GoogleInstanceService,
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

	// Atach the Disk to the VM
	deviceName, err := ad.vmService.AttachDisk(string(vmCID), disk.SelfLink)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	// Read VM agent settings
	agentSettings, err := ad.registryClient.Fetch(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	// Update VM agent settings
	newAgentSettings := agentSettings.AttachPersistentDisk(string(diskCID), deviceName)
	err = ad.registryClient.Update(string(vmCID), newAgentSettings)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	return nil, nil
}
