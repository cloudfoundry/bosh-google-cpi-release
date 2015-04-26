package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	"github.com/frodenas/bosh-google-cpi/google/address"
	"github.com/frodenas/bosh-google-cpi/google/client"
	"github.com/frodenas/bosh-google-cpi/google/disk"
	"github.com/frodenas/bosh-google-cpi/google/disk_type"
	"github.com/frodenas/bosh-google-cpi/google/image"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/google/machine_type"
	"github.com/frodenas/bosh-google-cpi/google/network"
	"github.com/frodenas/bosh-google-cpi/google/operation"
	"github.com/frodenas/bosh-google-cpi/google/snapshot"
	"github.com/frodenas/bosh-google-cpi/google/target_pool"

	"github.com/frodenas/bosh-google-cpi/registry"
)

type concreteFactory struct {
	availableActions map[string]Action
}

func NewConcreteFactory(
	googleClient gclient.GoogleClient,
	uuidGen boshuuid.Generator,
	options ConcreteFactoryOptions,
	logger boshlog.Logger,
) concreteFactory {
	operationService := goperation.NewGoogleOperationService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	addressService := gaddress.NewGoogleAddressService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	diskService := gdisk.NewGoogleDiskService(
		googleClient.Project(),
		googleClient.ComputeService(),
		operationService,
		uuidGen,
		logger,
	)

	diskTypeService := gdisktype.NewGoogleDiskTypeService(
		googleClient.Project(),
		googleClient.ComputeService(),
		logger,
	)

	machineTypeService := gmachinetype.NewGoogleMachineTypeService(
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

	registryService := registry.NewRegistryService(
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

	return concreteFactory{
		availableActions: map[string]Action{
			// Disk management
			"create_disk": NewCreateDisk(diskService, diskTypeService, vmService, googleClient.DefaultZone()),
			"delete_disk": NewDeleteDisk(diskService),
			"attach_disk": NewAttachDisk(diskService, vmService, registryService),
			"detach_disk": NewDetachDisk(diskService, vmService, registryService),

			// Snapshot management
			"snapshot_disk":   NewSnapshotDisk(snapshotService, diskService),
			"delete_snapshot": NewDeleteSnapshot(snapshotService),

			// Stemcell management
			"create_stemcell": NewCreateStemcell(stemcellService),
			"delete_stemcell": NewDeleteStemcell(stemcellService),

			// VM management
			"create_vm":          NewCreateVM(vmService, addressService, diskService, machineTypeService, networkService, stemcellService, targetPoolService, registryService, options.Agent, googleClient.DefaultZone()),
			"configure_networks": NewConfigureNetworks(vmService, addressService, networkService, targetPoolService, registryService),
			"delete_vm":          NewDeleteVM(vmService, registryService),
			"reboot_vm":          NewRebootVM(vmService),
			"set_vm_metadata":    NewSetVMMetadata(vmService),
			"has_vm":             NewHasVM(vmService),
			"get_disks":          NewGetDisks(vmService),
		},
	}
}

func (f concreteFactory) Create(method string) (Action, error) {
	action, found := f.availableActions[method]
	if !found {
		return nil, bosherr.Errorf("Could not create action with method %s", method)
	}

	return action, nil
}
