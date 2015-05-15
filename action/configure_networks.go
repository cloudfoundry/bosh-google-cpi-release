package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/address_service"
	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	"github.com/frodenas/bosh-google-cpi/google/network_service"
	"github.com/frodenas/bosh-google-cpi/google/target_pool_service"

	"github.com/frodenas/bosh-registry/client"
)

type ConfigureNetworks struct {
	vmService         ginstance.InstanceService
	addressService    address.Service
	networkService    gnetwork.NetworkService
	targetPoolService gtargetpool.TargetPoolService
	registryClient    registry.Client
}

func NewConfigureNetworks(
	vmService ginstance.InstanceService,
	addressService address.Service,
	networkService gnetwork.NetworkService,
	targetPoolService gtargetpool.TargetPoolService,
	registryClient registry.Client,
) ConfigureNetworks {
	return ConfigureNetworks{
		vmService:         vmService,
		addressService:    addressService,
		networkService:    networkService,
		targetPoolService: targetPoolService,
		registryClient:    registryClient,
	}
}

func (rv ConfigureNetworks) Run(vmCID VMCID, networks Networks) (interface{}, error) {
	// Parse networks
	vmNetworks := networks.AsInstanceServiceNetworks()
	instanceNetworks := ginstance.NewGoogleInstanceNetworks(vmNetworks, rv.addressService, rv.networkService, rv.targetPoolService)
	if err := instanceNetworks.Validate(); err != nil {
		return "", bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Update networks
	err := rv.vmService.UpdateNetworkConfiguration(string(vmCID), instanceNetworks)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Read VM agent settings
	agentSettings, err := rv.registryClient.Fetch(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Update VM agent settings
	agentNetworks := networks.AsRegistryNetworks()
	newAgentSettings := agentSettings.ConfigureNetworks(agentNetworks)
	err = rv.registryClient.Update(string(vmCID), newAgentSettings)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	return nil, nil
}
