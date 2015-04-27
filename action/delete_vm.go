package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/address"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/google/network"
	"github.com/frodenas/bosh-google-cpi/google/target_pool"
	"github.com/frodenas/bosh-google-cpi/registry"
)

type DeleteVM struct {
	vmService         ginstance.GoogleInstanceService
	addressService    gaddress.GoogleAddressService
	networkService    gnetwork.GoogleNetworkService
	targetPoolService gtargetpool.GoogleTargetPoolService
	registryService   registry.RegistryService
}

func NewDeleteVM(
	vmService ginstance.GoogleInstanceService,
	addressService gaddress.GoogleAddressService,
	networkService gnetwork.GoogleNetworkService,
	targetPoolService gtargetpool.GoogleTargetPoolService,
	registryService registry.RegistryService,
) DeleteVM {
	return DeleteVM{
		vmService:         vmService,
		addressService:    addressService,
		networkService:    networkService,
		targetPoolService: targetPoolService,
		registryService:   registryService,
	}
}

func (dv DeleteVM) Run(vmCID VMCID) (interface{}, error) {
	// Delete VM networks
	var networks Networks
	vmNetworks := networks.AsGoogleInstanceNetworks()
	instanceNetworks := ginstance.NewGoogleInstanceNetworks(vmNetworks, dv.addressService, dv.networkService, dv.targetPoolService)

	err := dv.vmService.DeleteNetworkConfiguration(string(vmCID), instanceNetworks)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM
	err = dv.vmService.Delete(string(vmCID))
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
