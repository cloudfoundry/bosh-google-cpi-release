package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/disk_type_service"
	"bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/util"
)

type CreateDisk struct {
	diskService     disk.Service
	diskTypeService disktype.Service
	vmService       instance.Service
	defaultZone     string
}

func NewCreateDisk(
	diskService disk.Service,
	diskTypeService disktype.Service,
	vmService instance.Service,
	defaultZone string,
) CreateDisk {
	return CreateDisk{
		diskService:     diskService,
		diskTypeService: diskTypeService,
		vmService:       vmService,
		defaultZone:     defaultZone,
	}
}

func (cd CreateDisk) Run(size int, cloudProps DiskCloudProperties, vmCID VMCID) (DiskCID, error) {
	var diskType string
	var zone = cd.defaultZone

	// Find the VM (if provided) so we can create the disk in the same zone
	if vmCID != "" {
		vm, found, err := cd.vmService.Find(string(vmCID), "")
		if err != nil {
			return "", bosherr.WrapError(err, "Creating disk")
		}
		if !found {
			return "", api.NewVMNotFoundError(string(vmCID))
		}

		zone = vm.Zone
	}

	// Find the Disk Type (if provided)
	if cloudProps.DiskType != "" {
		dt, found, err := cd.diskTypeService.Find(cloudProps.DiskType, zone)
		if err != nil {
			return "", bosherr.WrapError(err, "Creating disk")
		}
		if !found {
			return "", bosherr.WrapErrorf(err, "Creating disk: Disk Type '%s' does not exists", cloudProps.DiskType)
		}

		diskType = dt.SelfLink
	}

	// Create the Disk
	disk, err := cd.diskService.Create(util.ConvertMib2Gib(size), diskType, zone)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating disk")
	}

	return DiskCID(disk), nil
}
