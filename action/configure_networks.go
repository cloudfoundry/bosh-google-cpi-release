package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/address"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/google/network"
	"github.com/frodenas/bosh-google-cpi/google/target_pool"
	"github.com/frodenas/bosh-google-cpi/registry"
)

type ConfigureNetworks struct {
	vmService         ginstance.GoogleInstanceService
	addressService    gaddress.GoogleAddressService
	networkService    gnetwork.GoogleNetworkService
	targetPoolService gtargetpool.GoogleTargetPoolService
	registryService   registry.RegistryService
}

func NewConfigureNetworks(
	vmService ginstance.GoogleInstanceService,
	addressService gaddress.GoogleAddressService,
	networkService gnetwork.GoogleNetworkService,
	targetPoolService gtargetpool.GoogleTargetPoolService,
	registryService registry.RegistryService,
) ConfigureNetworks {
	return ConfigureNetworks{
		vmService:         vmService,
		addressService:    addressService,
		networkService:    networkService,
		targetPoolService: targetPoolService,
		registryService:   registryService,
	}
}

func (rv ConfigureNetworks) Run(vmCID VMCID, networks Networks) (interface{}, error) {
	// Parse networks
	vmNetworks := networks.AsGoogleInstanceNetworks()
	instanceNetworks := ginstance.NewGoogleInstanceNetworks(vmNetworks, rv.addressService, rv.networkService, rv.targetPoolService)
	if err := instanceNetworks.Validate(); err != nil {
		return "", bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Update networks
	err := rv.vmService.UpdateNetworkConfiguration(string(vmCID), instanceNetworks)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Read VM agent settings
	agentSettings, err := rv.registryService.Fetch(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Update VM agent settings
	agentNetworks := networks.AsAgentNetworks()
	newAgentSettings := agentSettings.ConfigureNetworks(agentNetworks)
	err = rv.registryService.Update(string(vmCID), newAgentSettings)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	return nil, nil
}
