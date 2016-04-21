package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/instance_service"

	"bosh-google-cpi/registry"
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
	if err := dv.vmService.DeleteNetworkConfiguration(string(vmCID)); err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM
	if err := dv.vmService.Delete(string(vmCID)); err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM agent settings
	if err := dv.registryClient.Delete(string(vmCID)); err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	return nil, nil
}
