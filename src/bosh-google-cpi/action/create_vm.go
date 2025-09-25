package action

import (
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/acceleratortype"
	"bosh-google-cpi/google/disk"
	"bosh-google-cpi/google/disktype"
	"bosh-google-cpi/google/image"
	"bosh-google-cpi/google/instance"
	"bosh-google-cpi/google/machinetype"
	"bosh-google-cpi/registry"
	"bosh-google-cpi/util"
)

type createVMBase struct {
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

type CreateVMV1 struct {
	createVMBase
}

type CreateVMV2 struct {
	createVMBase
}

func NewCreateVMV1(
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
) CreateVMV1 {
	return CreateVMV1{
		createVMBase{
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
		},
	}
}

func NewCreateVMV2(
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
) CreateVMV2 {
	return CreateVMV2{
		createVMBase{
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
		},
	}
}

func (cv1 CreateVMV1) Run(agentID string, stemcellCID StemcellCID, cloudProps VMCloudProperties, networks Networks, disks []DiskCID, env Environment) (interface{}, error) {
	vmCid, _, err := cv1.createVMBase.Run(agentID, stemcellCID, cloudProps, networks, disks, env)
	if err != nil {
		return nil, err
	}

	return vmCid, nil
}

func (cv2 CreateVMV2) Run(agentID string, stemcellCID StemcellCID, cloudProps VMCloudProperties, networks Networks, disks []DiskCID, env Environment) (interface{}, error) {
	vmCid, networks, err := cv2.createVMBase.Run(agentID, stemcellCID, cloudProps, networks, disks, env)
	if err != nil {
		return nil, err
	}

	return []interface{}{vmCid, networks}, nil
}

func (cv createVMBase) Run(agentID string, stemcellCID StemcellCID, cloudProps VMCloudProperties, networks Networks, disks []DiskCID, env Environment) (VMCID, Networks, error) {
	// Find zone
	zone, err := cv.findZone(cloudProps.Zone, disks)
	if err != nil {
		return "", nil, err
	}

	// Find stemcell
	stemcellLink, err := cv.findStemcellLink(string(stemcellCID))
	if err != nil {
		return "", nil, err
	}

	// Find machine type
	machineTypeLink, err := cv.findMachineTypeLink(cloudProps, zone)
	if err != nil {
		return "", nil, err
	}

	// Find the root Disk Type
	rootDiskTypeLink, err := cv.findRootDiskTypeLink(cloudProps.RootDiskType, zone)
	if err != nil {
		return "", nil, err
	}

	// Find Accelerator Type
	acceleratorTypeLinks, err := cv.findAcceleratorTypeLinks(cloudProps.Accelerators, zone)
	if err != nil {
		return "", nil, err
	}

	// Parse networks
	vmNetworks := networks.AsInstanceServiceNetworks()
	if err = vmNetworks.Validate(); err != nil {
		return "", nil, bosherr.WrapError(err, "Creating VM")
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
				safeTag, _ := instance.SafeLabel(tag.(string), instance.LabelKey) //nolint:errcheck
				cloudProps.Tags = append(cloudProps.Tags, safeTag)
			}
		}
	}

	// Validate VM tags and labels
	if err = cloudProps.Validate(); err != nil {
		return "", nil, bosherr.WrapError(err, "Creating VM")
	}

	bs, err := parseBackendService(cloudProps.BackendService)
	if err != nil {
		return "", nil, bosherr.WrapErrorf(err, "Parsing BackendService %#v", cloudProps.BackendService)
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
		NodeGroup:         cloudProps.NodeGroup,
		EphemeralDiskType: ephemeral,
	}

	// Create VM
	vm, err := cv.vmService.Create(vmProps, vmNetworks, cv.registryOptions.EndpointWithCredentials())
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", nil, err
		}
		return "", nil, bosherr.WrapError(err, "Creating VM")
	}

	// If any of the below code fails, we must delete the created vm
	defer func() {
		if err != nil {
			cv.vmService.CleanUp(vm)
		}
	}()

	// Create VM settings
	agentNetworks := networks.AsRegistryNetworks()
	agentSettings := registry.NewAgentSettings(agentID, vm, agentNetworks, registry.EnvSettings(env), cv.agentOptions)

	if ephemeral == "local-ssd" {
		agentSettings.Disks.Ephemeral = "/dev/nvme0n1"
	}

	if err = cv.registryClient.Update(vm, agentSettings); err != nil {
		return "", nil, bosherr.WrapErrorf(err, "Creating VM")
	}

	return VMCID(vm), networks, nil
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

func (cv createVMBase) findZone(zoneName string, disks []DiskCID) (string, error) {
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

	return "", fmt.Errorf("Could not find zone %q", zoneName) //nolint:staticcheck
}

func isGcpImageURL(s string) bool {
	return strings.HasPrefix(s, "https://www.googleapis.com/compute/v1/projects/")
}

func (cv createVMBase) findStemcellLink(stemcellID string) (string, error) {
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

func (cv createVMBase) findMachineTypeLink(cloudProps VMCloudProperties, zone string) (string, error) {
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

		machineTypeLink = cv.machineTypeService.CustomLink(cloudProps.CPU, cloudProps.RAM, zone, cloudProps.MachineSeries)
	}

	return machineTypeLink, nil
}

func (cv createVMBase) findRootDiskSizeGb(rootDiskSizeGb int) int {
	diskSizeGb := cv.defaultRootDiskSizeGb
	if rootDiskSizeGb > 0 {
		diskSizeGb = rootDiskSizeGb
	}

	return diskSizeGb
}

func (cv createVMBase) findRootDiskTypeLink(diskTypeName string, zone string) (string, error) {
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
func (cv createVMBase) findAcceleratorTypeLinks(accelerators []Accelerator, zone string) ([]instance.Accelerator, error) {
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
