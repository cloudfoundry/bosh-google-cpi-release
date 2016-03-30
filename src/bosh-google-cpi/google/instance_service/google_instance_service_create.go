package instance

import (
	"encoding/json"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

const defaultRootDiskSizeGb = 10
const userDataKey = "user_data"

func (i GoogleInstanceService) Create(vmProps *Properties, networks Networks, registryEndpoint string) (string, error) {
	uuidStr, err := i.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random Google Instance name")
	}

	instanceName := fmt.Sprintf("%s-%s", googleInstanceNamePrefix, uuidStr)
	canIPForward := networks.CanIPForward()
	diskParams := i.createDiskParams(vmProps.Stemcell, vmProps.RootDiskSizeGb, vmProps.RootDiskType)
	metadataParams, err := i.createMatadataParams(instanceName, registryEndpoint, networks)
	if err != nil {
		return "", err
	}
	networkInterfacesParams, err := i.createNetworkInterfacesParams(networks)
	if err != nil {
		return "", err
	}
	schedulingParams := i.createSchedulingParams(vmProps.AutomaticRestart, vmProps.OnHostMaintenance, vmProps.Preemptible)
	serviceAccountsParams := i.createServiceAccountsParams(vmProps.ServiceScopes)
	tagsParams, err := i.createTagsParams(networks)
	if err != nil {
		return "", err
	}

	vm := &compute.Instance{
		Name:              instanceName,
		Description:       googleInstanceDescription,
		CanIpForward:      canIPForward,
		Disks:             diskParams,
		MachineType:       vmProps.MachineType,
		Metadata:          metadataParams,
		NetworkInterfaces: networkInterfacesParams,
		Scheduling:        schedulingParams,
		ServiceAccounts:   serviceAccountsParams,
		Tags:              tagsParams,
	}
	i.logger.Debug(googleInstanceServiceLogTag, "Creating Google Instance with params: %#v", vm)
	operation, err := i.computeService.Instances.Insert(i.project, util.ResourceSplitter(vmProps.Zone), vm).Do()
	if err != nil {
		i.logger.Debug(googleInstanceServiceLogTag, "Failed to create Google Instance: %#v", err)
		return "", api.NewVMCreationFailedError(true)
	}

	if _, err = i.operationService.Waiter(operation, vmProps.Zone, ""); err != nil {
		i.logger.Debug(googleInstanceServiceLogTag, "Failed to create Google Instance: %#v", err)
		i.CleanUp(vm.Name)
		return "", api.NewVMCreationFailedError(true)
	}

	return vm.Name, nil
}

func (i GoogleInstanceService) CleanUp(id string) {
	if err := i.Delete(id); err != nil {
		i.logger.Debug(googleInstanceServiceLogTag, "Failed cleaning up Google Instance '%s': %#v", id, err)
	}
}

func (i GoogleInstanceService) createDiskParams(stemcell string, diskSize int, diskType string) []*compute.AttachedDisk {
	var disks []*compute.AttachedDisk

	if diskSize == 0 {
		diskSize = defaultRootDiskSizeGb
	}
	disk := &compute.AttachedDisk{
		AutoDelete: true,
		Boot:       true,
		InitializeParams: &compute.AttachedDiskInitializeParams{
			DiskSizeGb:  int64(diskSize),
			DiskType:    diskType,
			SourceImage: stemcell,
		},
		Mode: "READ_WRITE",
		Type: "PERSISTENT",
	}
	disks = append(disks, disk)

	return disks
}

func (i GoogleInstanceService) createMatadataParams(name string, regEndpoint string, networks Networks) (*compute.Metadata, error) {
	serverName := GoogleUserDataServerName{Name: name}
	registryEndpoint := GoogleUserDataRegistryEndpoint{Endpoint: regEndpoint}
	userData := GoogleUserData{Server: serverName, Registry: registryEndpoint}

	if networkDNS := networks.DNS(); len(networkDNS) > 0 {
		userData.DNS = GoogleUserDataDNSItems{NameServer: networkDNS}
	}

	ud, err := json.Marshal(userData)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Marshalling user data")
	}

	var metadataItems []*compute.MetadataItems
	userDataValue := string(ud)
	metadataItem := &compute.MetadataItems{Key: userDataKey, Value: &userDataValue}
	metadataItems = append(metadataItems, metadataItem)
	metadata := &compute.Metadata{Items: metadataItems}

	return metadata, nil
}

func (i GoogleInstanceService) createNetworkInterfacesParams(networks Networks) ([]*compute.NetworkInterface, error) {
	network, found, err := i.networkService.Find(networks.NetworkName())
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, bosherr.WrapErrorf(err, "Network '%s' does not exists", networks.NetworkName())
	}

	subnetworkLink := ""
	if networks.SubnetworkName() != "" {
		subnetwork, found, err := i.subnetworkService.Find(networks.SubnetworkName(), "")
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, bosherr.WrapErrorf(err, "Subnetwork '%s' does not exists", networks.SubnetworkName())
		}
		subnetworkLink = subnetwork.SelfLink
	}

	var networkInterfaces []*compute.NetworkInterface
	var accessConfigs []*compute.AccessConfig

	vipNetwork := networks.VipNetwork()
	if networks.EphemeralExternalIP() || vipNetwork.IP != "" {
		accessConfig := &compute.AccessConfig{
			Name: "External NAT",
			Type: "ONE_TO_ONE_NAT",
		}
		if vipNetwork.IP != "" {
			accessConfig.NatIP = vipNetwork.IP
		}
		accessConfigs = append(accessConfigs, accessConfig)
	}

	networkInterface := &compute.NetworkInterface{
		Network:       network.SelfLink,
		Subnetwork:    subnetworkLink,
		AccessConfigs: accessConfigs,
	}
	networkInterfaces = append(networkInterfaces, networkInterface)

	return networkInterfaces, nil
}

func (i GoogleInstanceService) createSchedulingParams(
	automaticRestart bool,
	onHostMaintenance string,
	preemptible bool,
) *compute.Scheduling {
	if preemptible {
		return &compute.Scheduling{Preemptible: preemptible}
	}

	scheduling := &compute.Scheduling{
		AutomaticRestart:  automaticRestart,
		OnHostMaintenance: onHostMaintenance,
		Preemptible:       preemptible,
	}

	if onHostMaintenance == "" {
		scheduling.OnHostMaintenance = "MIGRATE"
	}

	return scheduling
}

func (i GoogleInstanceService) createServiceAccountsParams(serviceScopes ServiceScopes) []*compute.ServiceAccount {
	var serviceAccounts []*compute.ServiceAccount

	if len(serviceScopes) > 0 {
		var scopes []string
		for _, serviceScope := range serviceScopes {
			scopes = append(scopes, fmt.Sprintf("https://www.googleapis.com/auth/%s", serviceScope))
		}
		serviceAccount := &compute.ServiceAccount{
			Email:  "default",
			Scopes: scopes,
		}
		serviceAccounts = append(serviceAccounts, serviceAccount)
	}

	return serviceAccounts
}

func (i GoogleInstanceService) createTagsParams(networks Networks) (*compute.Tags, error) {
	tags := &compute.Tags{}

	for _, tag := range networks.Tags() {
		tags.Items = append(tags.Items, tag)
	}

	return tags, nil
}
