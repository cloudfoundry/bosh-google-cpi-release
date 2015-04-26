package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/disk"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/registry"
)

type DetachDisk struct {
	diskService     gdisk.GoogleDiskService
	vmService       ginstance.GoogleInstanceService
	registryService registry.RegistryService
}

func NewDetachDisk(
	diskService gdisk.GoogleDiskService,
	vmService ginstance.GoogleInstanceService,
	registryService registry.RegistryService,
) DetachDisk {
	return DetachDisk{
		diskService:     diskService,
		vmService:       vmService,
		registryService: registryService,
	}
}

func (dd DetachDisk) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	// Detach the disk
	err := dd.vmService.DetachDisk(string(vmCID), string(diskCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Detaching disk '%s' from vm '%s", diskCID, vmCID)
	}

	// Read VM agent settings
	agentSettings, err := dd.registryService.Fetch(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Detaching disk '%s' from vm '%s", diskCID, vmCID)
	}

	// Update VM agent settings
	newAgentSettings := agentSettings.DetachPersistentDisk(string(diskCID))
	err = dd.registryService.Update(string(vmCID), newAgentSettings)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Detaching disk '%s' from vm '%s", diskCID, vmCID)
	}

	return nil, nil
}
