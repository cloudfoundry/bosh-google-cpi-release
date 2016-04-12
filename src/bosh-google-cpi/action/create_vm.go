package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/disk_type_service"
	"bosh-google-cpi/google/image_service"
	"bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/google/machine_type_service"
	"bosh-google-cpi/util"

	"bosh-google-cpi/registry"
)

type CreateVM struct {
	vmService             instance.Service
	diskService           disk.Service
	diskTypeService       disktype.Service
	imageService          image.Service
	machineTypeService    machinetype.Service
	registryClient        registry.Client
	registryOptions       registry.ClientOptions
	agentOptions          registry.AgentOptions
	defaultRootDiskSizeGb int
	defaultRootDiskType   string
	defaultZone           string
}

func NewCreateVM(
	vmService instance.Service,
	diskService disk.Service,
	diskTypeService disktype.Service,
	imageService image.Service,
	machineTypeService machinetype.Service,
	registryClient registry.Client,
	registryOptions registry.ClientOptions,
	agentOptions registry.AgentOptions,
	defaultRootDiskSizeGb int,
	defaultRootDiskType string,
	defaultZone string,
) CreateVM {
	return CreateVM{
		vmService:             vmService,
		diskService:           diskService,
		diskTypeService:       diskTypeService,
		imageService:          imageService,
		machineTypeService:    machineTypeService,
		registryClient:        registryClient,
		registryOptions:       registryOptions,
		agentOptions:          agentOptions,
		defaultRootDiskSizeGb: defaultRootDiskSizeGb,
		defaultRootDiskType:   defaultRootDiskType,
		defaultZone:           defaultZone,
	}
}

func (cv CreateVM) Run(agentID string, stemcellCID StemcellCID, cloudProps VMCloudProperties, networks Networks, disks []DiskCID, env Environment) (VMCID, error) {
	// Find zone
	zone, err := cv.findZone(cloudProps.Zone, disks)
	if err != nil {
		return "", err
	}

	// Find stemcell
	stemcellLink, err := cv.findStemcellLink(string(stemcellCID))
	if err != nil {
		return "", err
	}

	// Find machine type
	machineTypeLink, err := cv.findMachineTypeLink(cloudProps, zone)
	if err != nil {
		return "", err
	}

	// Find the root Disk Type
	rootDiskTypeLink, err := cv.findRootDiskTypeLink(cloudProps.RootDiskType, zone)
	if err != nil {
		return "", err
	}

	// Parse networks
	vmNetworks := networks.AsInstanceServiceNetworks()
	if err = vmNetworks.Validate(); err != nil {
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Parse VM properties
	vmProps := &instance.Properties{
		Zone:              zone,
		Stemcell:          stemcellLink,
		MachineType:       machineTypeLink,
		RootDiskSizeGb:    cv.findRootDiskSizeGb(cloudProps.RootDiskSizeGb),
		RootDiskType:      rootDiskTypeLink,
		AutomaticRestart:  cloudProps.AutomaticRestart,
		OnHostMaintenance: cloudProps.OnHostMaintenance,
		Preemptible:       cloudProps.Preemptible,
		ServiceScopes:     instance.ServiceScopes(cloudProps.ServiceScopes),
	}

	// Create VM
	vm, err := cv.vmService.Create(vmProps, vmNetworks, cv.registryOptions.EndpointWithCredentials())
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
	if err = cv.vmService.AddNetworkConfiguration(vm, vmNetworks); err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", err
		}
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Create VM settings
	agentNetworks := networks.AsRegistryNetworks()
	agentSettings := registry.NewAgentSettings(agentID, vm, agentNetworks, registry.EnvSettings(env), cv.agentOptions)
	if err = cv.registryClient.Update(vm, agentSettings); err != nil {
		return "", bosherr.WrapErrorf(err, "Creating VM")
	}

	return VMCID(vm), nil
}

func (cv CreateVM) findZone(zoneName string, disks []DiskCID) (string, error) {
	zones := make(map[string]struct{})
	if zoneName != "" {
		zones[zoneName] = struct{}{}
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

	for zone := range zones {
		return zone, nil
	}

	return cv.defaultZone, nil
}

func (cv CreateVM) findStemcellLink(stemcellID string) (string, error) {
	stemcell, found, err := cv.imageService.Find(stemcellID)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating vm")
	}
	if !found {
		return "", bosherr.WrapErrorf(err, "Creating vm: Stemcell '%s' does not exists", stemcellID)
	}

	return stemcell.SelfLink, nil
}

func (cv CreateVM) findMachineTypeLink(cloudProps VMCloudProperties, zone string) (string, error) {
	machineTypeLink := ""
	if cloudProps.MachineType != "" {
		if cloudProps.CPU != 0 || cloudProps.RAM != 0 {
			return "", bosherr.Error("Creating vm: 'machine_type' and 'cpu' or 'ram' cannot be provided together")
		}

		machineType, found, err := cv.machineTypeService.Find(cloudProps.MachineType, zone)
		if err != nil {
			return "", bosherr.WrapError(err, "Creating vm")
		}
		if !found {
			return "", bosherr.WrapErrorf(err, "Creating vm: Machine Type '%s' does not exists", cloudProps.MachineType)
		}
		machineTypeLink = machineType.SelfLink
	} else {
		if cloudProps.CPU == 0 || cloudProps.RAM == 0 {
			return "", bosherr.Error("Creating vm: 'machine_type' or 'cpu' and 'ram' must be provided")
		}

		machineTypeLink = cv.machineTypeService.CustomLink(cloudProps.CPU, cloudProps.RAM, zone)
	}

	return machineTypeLink, nil
}

func (cv CreateVM) findRootDiskSizeGb(rootDiskSizeGb int) int {
	diskSizeGb := cv.defaultRootDiskSizeGb
	if rootDiskSizeGb > 0 {
		diskSizeGb = rootDiskSizeGb
	}

	return diskSizeGb
}

func (cv CreateVM) findRootDiskTypeLink(diskTypeName string, zone string) (string, error) {
	diskType := cv.defaultRootDiskType
	if diskTypeName != "" {
		diskType = diskTypeName
	}

	if diskType != "" {
		dt, found, err := cv.diskTypeService.Find(diskType, zone)
		if err != nil {
			return "", bosherr.WrapError(err, "Creating vm")
		}
		if !found {
			return "", bosherr.WrapErrorf(err, "Creating vm: Root Disk Type '%s' does not exists", diskTypeName)
		}

		return dt.SelfLink, nil
	}

	return "", nil
}
