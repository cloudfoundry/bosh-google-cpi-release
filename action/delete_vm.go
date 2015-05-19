package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/instance_service"

	"github.com/frodenas/bosh-registry/client"
)

type DeleteVM struct {
	vmService      instance.Service
	registryClient registry.Client
}

func NewDeleteVM(
	vmService instance.Service,
	registryClient registry.Client,
) DeleteVM {
	return DeleteVM{
		vmService:      vmService,
		registryClient: registryClient,
	}
}

func (dv DeleteVM) Run(vmCID VMCID) (interface{}, error) {
	// Delete VM networks
	err := dv.vmService.DeleteNetworkConfiguration(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM
	err = dv.vmService.Delete(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM agent settings
	err = dv.registryClient.Delete(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	return nil, nil
}
