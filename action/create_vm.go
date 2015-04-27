package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/address"
	"github.com/frodenas/bosh-google-cpi/google/disk"
	"github.com/frodenas/bosh-google-cpi/google/image"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/google/machine_type"
	"github.com/frodenas/bosh-google-cpi/google/network"
	"github.com/frodenas/bosh-google-cpi/google/target_pool"
	"github.com/frodenas/bosh-google-cpi/google/util"
	"github.com/frodenas/bosh-google-cpi/registry"
)

type CreateVM struct {
	vmService          ginstance.GoogleInstanceService
	addressService     gaddress.GoogleAddressService
	diskService        gdisk.GoogleDiskService
	machineTypeService gmachinetype.GoogleMachineTypeService
	networkService     gnetwork.GoogleNetworkService
	stemcellService    gimage.GoogleImageService
	targetPoolService  gtargetpool.GoogleTargetPoolService
	registryService    registry.RegistryService
	agentOptions       registry.AgentOptions
	defaultZone        string
}

func NewCreateVM(
	vmService ginstance.GoogleInstanceService,
	addressService gaddress.GoogleAddressService,
	diskService gdisk.GoogleDiskService,
	machineTypeService gmachinetype.GoogleMachineTypeService,
	networkService gnetwork.GoogleNetworkService,
	stemcellService gimage.GoogleImageService,
	targetPoolService gtargetpool.GoogleTargetPoolService,
	registryService registry.RegistryService,
	agentOptions registry.AgentOptions,
	defaultZone string,
) CreateVM {
	return CreateVM{
		vmService:          vmService,
		addressService:     addressService,
		diskService:        diskService,
		machineTypeService: machineTypeService,
		networkService:     networkService,
		stemcellService:    stemcellService,
		targetPoolService:  targetPoolService,
		registryService:    registryService,
		agentOptions:       agentOptions,
		defaultZone:        defaultZone,
	}
}

func (cv CreateVM) Run(agentID string, stemcellCID StemcellCID, cloudProps VMCloudProperties, networks Networks, disks []DiskCID, env Environment) (VMCID, error) {
	// Find all affinity zones
	zones := make(map[string]struct{})
	if cloudProps.Zone != "" {
		zones[cloudProps.Zone] = struct{}{}
	}
	for _, diskCID := range disks {
		disk, found, err := cv.diskService.Find(string(diskCID), "")
		if err != nil {
			return "", bosherr.WrapError(err, "Creating vm")
		}
		if !found {
			return "", api.NewDiskNotFoundError(string(diskCID), false)
		}
		zones[gutil.ResourceSplitter(disk.Zone)] = struct{}{}
	}
	if len(zones) > 1 {
		return "", bosherr.Errorf("Creating vm: can't use multiple zones: '%v'", zones)
	}

	// Determine zone
	zone := cv.defaultZone
	for k := range zones {
		zone = k
		break
	}

	// Find stemcell
	stemcell, found, err := cv.stemcellService.Find(string(stemcellCID))
	if err != nil {
		return "", bosherr.WrapError(err, "Creating vm")
	}
	if !found {
		return "", bosherr.WrapErrorf(err, "Creating vm: Stemcell '%s' does not exists", stemcellCID)
	}

	// Find machine type
	if cloudProps.MachineType == "" {
		return "", bosherr.WrapError(err, "Creating vm: 'machine_type' must be provided")
	}
	machineType, found, err := cv.machineTypeService.Find(string(cloudProps.MachineType), zone)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating vm")
	}
	if !found {
		return "", bosherr.WrapErrorf(err, "Creating vm: Machine Type '%s' does not exists", cloudProps.MachineType)
	}

	// Parse networks
	vmNetworks := networks.AsGoogleInstanceNetworks()
	instanceNetworks := ginstance.NewGoogleInstanceNetworks(vmNetworks, cv.addressService, cv.networkService, cv.targetPoolService)
	if err = instanceNetworks.Validate(); err != nil {
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Parse VM properties
	vmProps := &ginstance.GoogleInstanceProperties{
		Zone:              zone,
		Stemcell:          stemcell.SelfLink,
		MachineType:       machineType.SelfLink,
		AutomaticRestart:  cloudProps.AutomaticRestart,
		OnHostMaintenance: cloudProps.OnHostMaintenance,
		ServiceScopes:     ginstance.GoogleInstanceServiceScopes(cloudProps.ServiceScopes),
	}

	// Create VM
	vm, err := cv.vmService.Create(vmProps, instanceNetworks, cv.registryService.PublicEndpoint())
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", err
		}
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Configure VM networks
	err = cv.vmService.AddNetworkConfiguration(vm, instanceNetworks)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", err
		}
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Create VM settings
	agentNetworks := networks.AsAgentNetworks()
	agentSettings := registry.NewAgentSettingsForVM(agentID, vm, agentNetworks, registry.EnvSettings(env), cv.agentOptions)
	err = cv.registryService.Update(vm, agentSettings)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Creating VM")
	}

	return VMCID(vm), nil
}
