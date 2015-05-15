package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/instance_service"

	"github.com/frodenas/bosh-registry/client"
)

type DetachDisk struct {
	vmService      ginstance.InstanceService
	registryClient registry.Client
}

func NewDetachDisk(
	vmService ginstance.InstanceService,
	registryClient registry.Client,
) DetachDisk {
	return DetachDisk{
		vmService:      vmService,
		registryClient: registryClient,
	}
}

func (dd DetachDisk) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	// Detach the disk
	err := dd.vmService.DetachDisk(string(vmCID), string(diskCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Detaching disk '%s' from vm '%s", diskCID, vmCID)
	}

	// Read VM agent settings
	agentSettings, err := dd.registryClient.Fetch(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Detaching disk '%s' from vm '%s", diskCID, vmCID)
	}

	// Update VM agent settings
	newAgentSettings := agentSettings.DetachPersistentDisk(string(diskCID))
	err = dd.registryClient.Update(string(vmCID), newAgentSettings)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Detaching disk '%s' from vm '%s", diskCID, vmCID)
	}

	return nil, nil
}
