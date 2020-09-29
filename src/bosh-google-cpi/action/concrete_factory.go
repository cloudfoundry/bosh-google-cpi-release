package action

import (
	"encoding/json"

	bogcconfig "bosh-google-cpi/google/config"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"bosh-google-cpi/config"
	"bosh-google-cpi/google/address_service"
	"bosh-google-cpi/google/backendservice_service"
	"bosh-google-cpi/google/client"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/disk_type_service"
	"bosh-google-cpi/google/image_service"
	"bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/google/machine_type_service"
	"bosh-google-cpi/google/network_service"
	"bosh-google-cpi/google/operation_service"
	"bosh-google-cpi/google/project_service"
	"bosh-google-cpi/google/snapshot_service"
	"bosh-google-cpi/google/subnetwork_service"
	"bosh-google-cpi/google/target_pool_service"

	"bosh-google-cpi/google/accelerator_type_service"
	"bosh-google-cpi/registry"
)

var GoogleClientFunc func(bogcconfig.Config, boshlog.Logger) (client.GoogleClient, error) = client.NewGoogleClient

type ConcreteFactory struct {
	uuidGen boshuuid.Generator
	cfg     config.Config
	logger  boshlog.Logger
}

func NewConcreteFactory(
	uuidGen boshuuid.Generator,
	cfg config.Config,
	logger boshlog.Logger,
) ConcreteFactory {
	return ConcreteFactory{uuidGen,
		cfg,
		logger}
}

func (f ConcreteFactory) Create(method string, ctx map[string]interface{}, apiVersion int) (Action, error) {
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Remarshaling")
	}

	err = json.Unmarshal(ctxBytes, &f.cfg.Cloud.Properties.Google)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Unmarshaling into google props")
	}

	googleClient, err := GoogleClientFunc(f.cfg.Cloud.Properties.Google, f.logger)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Building goog client")
	}

	operationService := operation.NewGoogleOperationService(
		googleClient.Project(),
		googleClient.ComputeService(),
		googleClient.ComputeBetaService(),
		f.logger,
	)

	addressService := address.NewGoogleAddressService(
		googleClient.Project(),
		googleClient.ComputeService(),
		f.logger,
	)

	diskService := disk.NewGoogleDiskService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		f.uuidGen,
		f.logger,
	)

	diskTypeService := disktype.NewGoogleDiskTypeService(
		googleClient.Project(),
		googleClient.ComputeService(),
		f.logger,
	)

	imageService := image.NewGoogleImageService(
		googleClient.Project(),
		googleClient.ComputeService(),
		googleClient.StorageService(),
		operationService,
		f.uuidGen,
		f.logger,
	)

	backendServiceService := backendservice.NewGoogleBackendServiceService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		f.logger,
	)

	machineTypeService := machinetype.NewGoogleMachineTypeService(
		googleClient.Project(),
		googleClient.ComputeService(),
		f.logger,
	)

	acceleratorTypeService := acceleratortype.NewGoogleAcceleratorTypeService(
		googleClient.Project(),
		googleClient.ComputeService(),
		f.logger,
	)

	projectService := project.NewGoogleProjectService(
		googleClient.Project(),
	)

	networkService := network.NewGoogleNetworkService(
		projectService,
		googleClient.ComputeService(),
		f.logger,
	)

	registryClient := registry.NewMetadataClient(
		googleClient,
		f.cfg.Cloud.Properties.Registry,
		f.logger,
	)

	snapshotService := snapshot.NewGoogleSnapshotService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		f.uuidGen,
		f.logger,
	)

	subnetworkService := subnetwork.NewGoogleSubnetworkService(
		projectService,
		googleClient.ComputeService(),
		f.logger,
	)

	targetPoolService := targetpool.NewGoogleTargetPoolService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		f.logger,
	)

	vmService := instance.NewGoogleInstanceService(
		googleClient.Project(),
		googleClient.ComputeService(),
		googleClient.ComputeBetaService(),
		addressService,
		backendServiceService,
		networkService,
		operationService,
		subnetworkService,
		targetPoolService,
		diskTypeService,
		f.uuidGen,
		f.logger,
	)

	actions := map[string]Action{
		// Disk management
		"create_disk": NewCreateDisk(
			diskService,
			diskTypeService,
			vmService,
		),
		"delete_disk": NewDeleteDisk(diskService),
		"attach_disk": f.selectAttachDisk(apiVersion, diskService, vmService, registryClient),
		"detach_disk": NewDetachDisk(vmService, registryClient),
		"has_disk":    NewHasDisk(diskService),

		// Snapshot management
		"snapshot_disk":   NewSnapshotDisk(snapshotService, diskService),
		"delete_snapshot": NewDeleteSnapshot(snapshotService),

		// Stemcell management
		"create_stemcell": NewCreateStemcell(imageService),
		"delete_stemcell": NewDeleteStemcell(imageService),

		// VM management
		"create_vm": f.selectCreateVM(
			apiVersion,
			vmService,
			diskService,
			diskTypeService,
			imageService,
			machineTypeService,
			acceleratorTypeService,
			registryClient,
			f.cfg.Cloud.Properties.Registry,
			f.cfg.Cloud.Properties.Agent,
			googleClient.DefaultRootDiskSizeGb(),
			googleClient.DefaultRootDiskType(),
		),
		"configure_networks": NewConfigureNetworks(vmService, registryClient),
		"delete_vm":          NewDeleteVM(vmService, registryClient),
		"reboot_vm":          NewRebootVM(vmService),
		"set_vm_metadata":    NewSetVMMetadata(vmService),
		"has_vm":             NewHasVM(vmService),
		"get_disks":          NewGetDisks(vmService),

		// Others:
		"info":                          NewInfo(),
		"ping":                          NewPing(),
		"calculate_vm_cloud_properties": NewCalculateVMCloudProperties(),

		// Not implemented:
		// current_vm_id
	}

	action, found := actions[method]
	if !found {
		return nil, bosherr.Errorf("Could not create action with method %s", method)
	}

	return action, nil
}

func (f ConcreteFactory) selectAttachDisk(
	apiVersion int,
	diskService disk.Service,
	vmService instance.Service,
	registryClient registry.Client,
) interface{} {
	if apiVersion == 2 {
		return NewAttachDiskV2(diskService, vmService, registryClient)
	} else {
		return NewAttachDiskV1(diskService, vmService, registryClient)
	}
}

func (f ConcreteFactory) selectCreateVM(
	apiVersion int,
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
) interface{} {
	if apiVersion == 2 {
		return NewCreateVMV2(
			vmService,
			diskService,
			diskTypeService,
			imageService,
			machineTypeService,
			acceleratorTypeService,
			registryClient,
			registryOptions,
			agentOptions,
			defaultRootDiskSizeGb,
			defaultRootDiskType,
		)
	} else {
		return NewCreateVMV1(
			vmService,
			diskService,
			diskTypeService,
			imageService,
			machineTypeService,
			acceleratorTypeService,
			registryClient,
			registryOptions,
			agentOptions,
			defaultRootDiskSizeGb,
			defaultRootDiskType,
		)
	}
}
