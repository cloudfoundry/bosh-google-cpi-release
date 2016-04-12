package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"bosh-google-cpi/google/address_service"
	"bosh-google-cpi/google/client"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/disk_type_service"
	"bosh-google-cpi/google/image_service"
	"bosh-google-cpi/google/instance_group_service"
	"bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/google/machine_type_service"
	"bosh-google-cpi/google/network_service"
	"bosh-google-cpi/google/operation_service"
	"bosh-google-cpi/google/snapshot_service"
	"bosh-google-cpi/google/subnetwork_service"
	"bosh-google-cpi/google/target_pool_service"

	"bosh-google-cpi/registry"
)

type ConcreteFactory struct {
	availableActions map[string]Action
}

func NewConcreteFactory(
	googleClient client.GoogleClient,
	uuidGen boshuuid.Generator,
	options ConcreteFactoryOptions,
	logger boshlog.Logger,
) ConcreteFactory {
	operationService := operation.NewGoogleOperationService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	addressService := address.NewGoogleAddressService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	diskService := disk.NewGoogleDiskService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		uuidGen,
		logger,
	)

	diskTypeService := disktype.NewGoogleDiskTypeService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	imageService := image.NewGoogleImageService(
		googleClient.Project(),
		googleClient.ComputeService(),
		googleClient.StorageService(),
		operationService,
		uuidGen,
		logger,
	)

	instanceGroupService := instancegroup.NewGoogleInstanceGroupService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		logger,
	)

	machineTypeService := machinetype.NewGoogleMachineTypeService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	networkService := network.NewGoogleNetworkService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	registryClient := registry.NewHTTPClient(
		options.Registry,
		logger,
	)

	snapshotService := snapshot.NewGoogleSnapshotService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		uuidGen,
		logger,
	)

	subnetworkService := subnetwork.NewGoogleSubnetworkService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	targetPoolService := targetpool.NewGoogleTargetPoolService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		logger,
	)

	vmService := instance.NewGoogleInstanceService(
		googleClient.Project(),
		googleClient.ComputeService(),
		addressService,
		instanceGroupService,
		networkService,
		operationService,
		subnetworkService,
		targetPoolService,
		uuidGen,
		logger,
	)

	return ConcreteFactory{
		availableActions: map[string]Action{
			// Disk management
			"create_disk": NewCreateDisk(
				diskService,
				diskTypeService,
				vmService,
				googleClient.DefaultZone(),
			),
			"delete_disk": NewDeleteDisk(diskService),
			"attach_disk": NewAttachDisk(diskService, vmService, registryClient),
			"detach_disk": NewDetachDisk(vmService, registryClient),

			// Snapshot management
			"snapshot_disk":   NewSnapshotDisk(snapshotService, diskService),
			"delete_snapshot": NewDeleteSnapshot(snapshotService),

			// Stemcell management
			"create_stemcell": NewCreateStemcell(imageService),
			"delete_stemcell": NewDeleteStemcell(imageService),

			// VM management
			"create_vm": NewCreateVM(
				vmService,
				diskService,
				diskTypeService,
				imageService,
				machineTypeService,
				registryClient,
				options.Registry,
				options.Agent,
				googleClient.DefaultRootDiskSizeGb(),
				googleClient.DefaultRootDiskType(),
				googleClient.DefaultZone(),
			),
			"configure_networks": NewConfigureNetworks(vmService, registryClient),
			"delete_vm":          NewDeleteVM(vmService, registryClient),
			"reboot_vm":          NewRebootVM(vmService),
			"set_vm_metadata":    NewSetVMMetadata(vmService),
			"has_vm":             NewHasVM(vmService),
			"get_disks":          NewGetDisks(vmService),

			// Others:
			"ping": NewPing(),

			// Not implemented:
			// current_vm_id
		},
	}
}

func (f ConcreteFactory) Create(method string) (Action, error) {
	action, found := f.availableActions[method]
	if !found {
		return nil, bosherr.Errorf("Could not create action with method %s", method)
	}

	return action, nil
}
