package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/address_service"
	"github.com/frodenas/bosh-google-cpi/google/disk_service"
	"github.com/frodenas/bosh-google-cpi/google/disk_type_service"
	"github.com/frodenas/bosh-google-cpi/google/image_service"
	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	"github.com/frodenas/bosh-google-cpi/google/machine_type_service"
	"github.com/frodenas/bosh-google-cpi/google/network_service"
	"github.com/frodenas/bosh-google-cpi/google/target_pool_service"
	"github.com/frodenas/bosh-google-cpi/util"

	"github.com/frodenas/bosh-registry/client"
)

type CreateVM struct {
	vmService          ginstance.InstanceService
	addressService     address.Service
	diskService        disk.Service
	diskTypeService    disktype.Service
	machineTypeService machinetype.Service
	networkService     gnetwork.NetworkService
	stemcellService    gimage.ImageService
	targetPoolService  gtargetpool.TargetPoolService
	registryClient     registry.Client
	registryOptions    registry.ClientOptions
	agentOptions       registry.AgentOptions
	defaultZone        string
}

func NewCreateVM(
	vmService ginstance.InstanceService,
	addressService address.Service,
	diskService disk.Service,
	diskTypeService disktype.Service,
	machineTypeService machinetype.Service,
	networkService gnetwork.NetworkService,
	stemcellService gimage.ImageService,
	targetPoolService gtargetpool.TargetPoolService,
	registryClient registry.Client,
	registryOptions registry.ClientOptions,
	agentOptions registry.AgentOptions,
	defaultZone string,
) CreateVM {
	return CreateVM{
		vmService:          vmService,
		addressService:     addressService,
		diskService:        diskService,
		diskTypeService:    diskTypeService,
		machineTypeService: machineTypeService,
		networkService:     networkService,
		stemcellService:    stemcellService,
		targetPoolService:  targetPoolService,
		registryClient:     registryClient,
		registryOptions:    registryOptions,
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
		zones[util.ResourceSplitter(disk.Zone)] = struct{}{}
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

	// Find the Disk Type (if provided)
	var diskType string
	if cloudProps.RootDiskType != "" {
		dt, found, err := cv.diskTypeService.Find(cloudProps.RootDiskType, zone)
		if err != nil {
			return "", bosherr.WrapError(err, "Creating vm")
		}
		if !found {
			return "", bosherr.WrapErrorf(err, "Creating vm: Root Disk Type '%s' does not exists", cloudProps.RootDiskType)
		}

		diskType = dt.SelfLink
	}

	// Parse networks
	vmNetworks := networks.AsInstanceServiceNetworks()
	instanceNetworks := ginstance.NewGoogleInstanceNetworks(vmNetworks, cv.addressService, cv.networkService, cv.targetPoolService)
	if err = instanceNetworks.Validate(); err != nil {
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Parse VM properties
	vmProps := &ginstance.InstanceProperties{
		Zone:              zone,
		Stemcell:          stemcell.SelfLink,
		MachineType:       machineType.SelfLink,
		RootDiskSizeGb:    cloudProps.RootDiskSizeGb,
		RootDiskType:      diskType,
		AutomaticRestart:  cloudProps.AutomaticRestart,
		OnHostMaintenance: cloudProps.OnHostMaintenance,
		ServiceScopes:     ginstance.InstanceServiceScopes(cloudProps.ServiceScopes),
	}

	// Create VM
	vm, err := cv.vmService.Create(vmProps, instanceNetworks, cv.registryOptions.Endpoint())
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", err
		}
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// If any of the below code fails, we must delete the created vm
	defer func() {
		if err != nil {
			cv.vmService.CleanUp(vm)
		}
	}()

	// Configure VM networks
	err = cv.vmService.AddNetworkConfiguration(vm, instanceNetworks)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", err
		}
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Create VM settings
	agentNetworks := networks.AsRegistryNetworks()
	agentSettings := registry.NewAgentSettings(agentID, vm, agentNetworks, registry.EnvSettings(env), cv.agentOptions)
	err = cv.registryClient.Update(vm, agentSettings)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Creating VM")
	}

	return VMCID(vm), nil
}
