package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/registry"
)

type DeleteVM struct {
	vmService       ginstance.GoogleInstanceService
	registryService registry.RegistryService
}

func NewDeleteVM(
	vmService ginstance.GoogleInstanceService,
	registryService registry.RegistryService,
) DeleteVM {
	return DeleteVM{
		vmService:       vmService,
		registryService: registryService,
	}
}

func (dv DeleteVM) Run(vmCID VMCID) (interface{}, error) {
	// Delete the VM
	err := dv.vmService.Delete(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM agent settings
	err = dv.registryService.Delete(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	return nil, nil
}
