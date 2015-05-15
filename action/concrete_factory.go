package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	"github.com/frodenas/bosh-google-cpi/google/address_service"
	"github.com/frodenas/bosh-google-cpi/google/client"
	"github.com/frodenas/bosh-google-cpi/google/disk_service"
	"github.com/frodenas/bosh-google-cpi/google/disk_type_service"
	"github.com/frodenas/bosh-google-cpi/google/image_service"
	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	"github.com/frodenas/bosh-google-cpi/google/machine_type_service"
	"github.com/frodenas/bosh-google-cpi/google/network_service"
	"github.com/frodenas/bosh-google-cpi/google/operation_service"
	"github.com/frodenas/bosh-google-cpi/google/snapshot_service"
	"github.com/frodenas/bosh-google-cpi/google/target_pool_service"

	"github.com/frodenas/bosh-registry/client"
)

type ConcreteFactory struct {
	availableActions map[string]Action
}

func NewConcreteFactory(
	googleClient gclient.GoogleClient,
	uuidGen boshuuid.Generator,
	options ConcreteFactoryOptions,
	logger boshlog.Logger,
) ConcreteFactory {
	operationService := goperation.NewGoogleOperationService(
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

	machineTypeService := machinetype.NewGoogleMachineTypeService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	networkService := gnetwork.NewGoogleNetworkService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	snapshotService := gsnapshot.NewGoogleSnapshotService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		uuidGen,
		logger,
	)

	registryClient := registry.NewHTTPClient(
		options.Registry,
		logger,
	)

	stemcellService := gimage.NewGoogleImageService(
		googleClient.Project(),
		googleClient.ComputeService(),
		googleClient.StorageService(),
		operationService,
		uuidGen,
		logger,
	)

	vmService := ginstance.NewGoogleInstanceService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		uuidGen,
		logger,
	)

	targetPoolService := gtargetpool.NewGoogleTargetPoolService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
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
			"create_stemcell": NewCreateStemcell(stemcellService),
			"delete_stemcell": NewDeleteStemcell(stemcellService),

			// VM management
			"create_vm": NewCreateVM(
				vmService,
				addressService,
				diskService,
				diskTypeService,
				machineTypeService,
				networkService,
				stemcellService,
				targetPoolService,
				registryClient,
				options.Registry,
				options.Agent,
				googleClient.DefaultZone(),
			),
			"configure_networks": NewConfigureNetworks(
				vmService,
				addressService,
				networkService,
				targetPoolService,
				registryClient,
			),
			"delete_vm": NewDeleteVM(
				vmService,
				addressService,
				networkService,
				targetPoolService,
				registryClient,
			),
			"reboot_vm":       NewRebootVM(vmService),
			"set_vm_metadata": NewSetVMMetadata(vmService),
			"has_vm":          NewHasVM(vmService),
			"get_disks":       NewGetDisks(vmService),

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
