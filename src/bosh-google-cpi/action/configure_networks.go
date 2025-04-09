package action

import (
	"bosh-google-cpi/api"
	instance "bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/registry"
)

type ConfigureNetworks struct {
	vmService      instance.Service
	registryClient registry.Client
}

func NewConfigureNetworks(
	vmService instance.Service,
	registryClient registry.Client,
) ConfigureNetworks {
	return ConfigureNetworks{
		vmService:      vmService,
		registryClient: registryClient,
	}
}

func (rv ConfigureNetworks) Run(vmCID VMCID, networks Networks) (interface{}, error) {
	return nil, api.NotSupportedError{}
}
