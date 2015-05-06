package ginstance

import (
	"encoding/json"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
)

const defaultRootDiskSizeGb = 10

func (i GoogleInstanceService) Create(vmProps *GoogleInstanceProperties, instanceNetworks GoogleInstanceNetworks, registryEndpoint string) (string, error) {
	uuidStr, err := i.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random Google Instance name")
	}

	instanceName := fmt.Sprintf("%s-%s", googleInstanceNamePrefix, uuidStr)
	canIPForward := instanceNetworks.CanIPForward()
	diskParams := i.createDiskParams(vmProps.Stemcell, vmProps.RootDiskSizeGb, vmProps.RootDiskType)
	metadataParams, err := i.createMatadataParams(instanceName, registryEndpoint, instanceNetworks)
	if err != nil {
		return "", err
	}
	networkInterfacesParams, err := instanceNetworks.NetworkInterfaces()
	if err != nil {
		return "", err
	}
	schedulingParams := i.createSchedulingParams(vmProps.AutomaticRestart, vmProps.OnHostMaintenance)
	serviceAccountsParams := i.createServiceAccountsParams(vmProps.ServiceScopes)
	tagsParams, err := instanceNetworks.Tags()
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
	operation, err := i.computeService.Instances.Insert(i.project, gutil.ResourceSplitter(vmProps.Zone), vm).Do()
	if err != nil {
		return "", api.NewVMCreationFailedError(true)
	}

	if _, err = i.operationService.Waiter(operation, vmProps.Zone, ""); err != nil {
		i.CleanUp(vm.Name)
		return "", api.NewVMCreationFailedError(true)
	}

	return vm.Name, nil
}

func (i GoogleInstanceService) CleanUp(id string) {
	err := i.Delete(id)
	if err != nil {
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

func (i GoogleInstanceService) createMatadataParams(name string, regEndpoint string, instanceNetworks GoogleInstanceNetworks) (*compute.Metadata, error) {
	serverName := GoogleUserDataServerName{Name: name}
	registryEndpoint := GoogleUserDataRegistryEndpoint{Endpoint: regEndpoint}
	userData := GoogleUserData{Server: serverName, Registry: registryEndpoint}

	if networkDNS := instanceNetworks.DNS(); len(networkDNS) > 0 {
		userData.DNS = GoogleUserDataDNSItems{NameServer: networkDNS}
	}

	ud, err := json.Marshal(userData)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Marshalling user data")
	}

	var metadataItems []*compute.MetadataItems
	metadataItem := &compute.MetadataItems{Key: "user_data", Value: string(ud)}
	metadataItems = append(metadataItems, metadataItem)
	metadata := &compute.Metadata{Items: metadataItems}

	return metadata, nil
}

func (i GoogleInstanceService) createSchedulingParams(automaticRestart bool, onHostMaintenance string) *compute.Scheduling {
	scheduling := &compute.Scheduling{
		AutomaticRestart:  automaticRestart,
		OnHostMaintenance: onHostMaintenance,
	}

	if onHostMaintenance == "" {
		scheduling.OnHostMaintenance = "MIGRATE"
	}

	return scheduling
}

func (i GoogleInstanceService) createServiceAccountsParams(serviceScopes GoogleInstanceServiceScopes) []*compute.ServiceAccount {
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
