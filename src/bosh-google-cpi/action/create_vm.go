package action

import (
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/disk_type_service"
	"bosh-google-cpi/google/image_service"
	"bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/google/machine_type_service"
	"bosh-google-cpi/util"

	"bosh-google-cpi/google/accelerator_type_service"
	"bosh-google-cpi/registry"
)

type CreateVM struct {
	vmService              instance.Service
	diskService            disk.Service
	diskTypeService        disktype.Service
	imageService           image.Service
	machineTypeService     machinetype.Service
	acceleratorTypeService acceleratortype.Service
	registryClient         registry.Client
	registryOptions        registry.ClientOptions
	agentOptions           registry.AgentOptions
	defaultRootDiskSizeGb  int
	defaultRootDiskType    string
}

func NewCreateVM(
	vmService instance.Service,
	diskService disk.Service,
	diskTypeService disktype.Service,
	imageService image.Service,
	machineTypeService machinetype.Service,
	acceleratorTypeService acceleratortype.Service,
	registryClient registry.Client,
	registryOptions registry.ClientOptions,
	agentOptions registry.AgentOptions,
	defaultRootDiskSizeGb int,
	defaultRootDiskType string,
) CreateVM {
	return CreateVM{
		vmService:              vmService,
		diskService:            diskService,
		diskTypeService:        diskTypeService,
		imageService:           imageService,
		machineTypeService:     machineTypeService,
		acceleratorTypeService: acceleratorTypeService,
		registryClient:         registryClient,
		registryOptions:        registryOptions,
		agentOptions:           agentOptions,
		defaultRootDiskSizeGb:  defaultRootDiskSizeGb,
		defaultRootDiskType:    defaultRootDiskType,
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

	// Find Accelerator Type
	acceleratorTypeLinks, err := cv.findAcceleratorTypeLinks(cloudProps.Accelerators, zone)
	if err != nil {
		return "", err
	}

	// Parse networks
	vmNetworks := networks.AsInstanceServiceNetworks()
	if err = vmNetworks.Validate(); err != nil {
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Certain properties defined in the Networks section of a manifest can be
	// overridden by VM properties. Here, we see if any of the VM properties
	// have been set and should override Network settings
	if cloudProps.IPForwarding != nil {
		vmNetworks.Network().IPForwarding = *cloudProps.IPForwarding
	}
	if cloudProps.EphemeralExternalIP != nil {
		vmNetworks.Network().EphemeralExternalIP = *cloudProps.EphemeralExternalIP
	}
	if vmNetworks.Network().Type == "dynamic" {
		vmNetworks.Network().IP = ""
	}

	// Extract any tags from env.bosh.groups
	if boshenv, ok := env["bosh"]; ok {
		if boshgroups, ok := boshenv.(map[string]interface{})["groups"]; ok {
			for _, tag := range boshgroups.([]interface{}) {
				// Ignore error as labels will be validated later
				safeTag, _ := instance.SafeLabel(tag.(string))
				cloudProps.Tags = append(cloudProps.Tags, safeTag)
			}
		}
	}

	// Validate VM tags and labels
	if err = cloudProps.Validate(); err != nil {
		return "", bosherr.WrapError(err, "Creating VM")
	}

	bs, err := parseBackendService(cloudProps.BackendService)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Parsing BackendService %#v", cloudProps.BackendService)
	}

	ephemeral := cloudProps.EphemeralDiskType

	// Parse VM properties
	vmProps := &instance.Properties{
		Zone:              zone,
		Name:              cloudProps.Name,
		Stemcell:          stemcellLink,
		MachineType:       machineTypeLink,
		RootDiskSizeGb:    cv.findRootDiskSizeGb(cloudProps.RootDiskSizeGb),
		RootDiskType:      rootDiskTypeLink,
		AutomaticRestart:  cloudProps.AutomaticRestart,
		OnHostMaintenance: cloudProps.OnHostMaintenance,
		Preemptible:       cloudProps.Preemptible,
		ServiceAccount:    instance.ServiceAccount(cloudProps.ServiceAccount),
		ServiceScopes:     instance.ServiceScopes(cloudProps.ServiceScopes),
		TargetPool:        cloudProps.TargetPool,
		BackendService:    bs,
		Tags:              cloudProps.Tags,
		Labels:            cloudProps.Labels,
		Accelerators:      acceleratorTypeLinks,
		EphemeralDiskType: ephemeral,
	}

	// Create VM
	vm, err, _ := cv.vmService.Create(vmProps, vmNetworks, cv.registryOptions.EndpointWithCredentials())
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

	if err != nil {
		return "", err
	}

	// Create VM settings
	agentNetworks := networks.AsRegistryNetworks()
	agentSettings := registry.NewAgentSettings(agentID, vm, agentNetworks, registry.EnvSettings(env), cv.agentOptions)

	if ephemeral == "local-ssd" {
		agentSettings.Disks.Ephemeral = "/dev/nvme0n1"
	}

	if err = cv.registryClient.Update(vm, agentSettings); err != nil {
		return "", bosherr.WrapErrorf(err, "Creating VM")
	}
	return VMCID(vm), nil
}

func extract(src map[string]interface{}, key string) string {
	if src[key] != nil {
		if val, ok := src[key].(string); ok {
			return val
		}
	}
	return ""
}

func parseBackendService(backendService interface{}) (instance.BackendService, error) {
	if backendService == nil {
		return instance.BackendService{}, nil
	}

	// backend_service: <name>
	if bs, ok := backendService.(string); ok {
		return instance.BackendService{Name: bs}, nil
	}

	//  backend_service:
	//    name: <name>
	//    scheme: <EXTERNAL|INTERNAL> (optional)
	bs := instance.BackendService{}
	if bsMap, ok := backendService.(map[string]string); ok {
		bs.Name = bsMap["name"]
	} else if bsMap, ok := backendService.(map[string]interface{}); ok {
		bs.Name = extract(bsMap, "name")
	} else if backendService != nil {
		return bs, bosherr.Errorf("unexpected type %T", backendService)
	}

	if bs.Name == "" {
		return bs, bosherr.Error("expected key: name")
	}

	return bs, nil
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

	return "", fmt.Errorf("Could not find zone %q", zoneName)
}

func isGcpImageURL(s string) bool {
	return strings.HasPrefix(s, "https://www.googleapis.com/compute/v1/projects/")
}

func (cv CreateVM) findStemcellLink(stemcellID string) (string, error) {
	if isGcpImageURL(stemcellID) {
		return stemcellID, nil
	}
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
func (cv CreateVM) findAcceleratorTypeLinks(accelerators []Accelerator, zone string) ([]instance.Accelerator, error) {
	if len(accelerators) == 0 {
		return nil, nil
	}

	acceleratorLinkTypes := []instance.Accelerator{}

	for _, acc := range accelerators {
		acceleratorType, found, err := cv.acceleratorTypeService.Find(acc.AcceleratorType, zone)
		if err != nil {
			return nil, bosherr.WrapError(err, "Creating vm")
		}
		if !found {
			return nil, bosherr.WrapErrorf(err, "Creating vm: Accelerator Type '%s' does not exists", acc.AcceleratorType)
		}
		updatedAcc := instance.Accelerator{
			AcceleratorType: acceleratorType.SelfLink,
			Count:           acc.Count,
		}
		acceleratorLinkTypes = append(acceleratorLinkTypes, updatedAcc)
	}

	return acceleratorLinkTypes, nil
}
