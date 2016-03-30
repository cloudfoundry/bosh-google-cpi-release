package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/instance_service"
)

type GetDisks struct {
	vmService instance.Service
}

func NewGetDisks(
	vmService instance.Service,
) GetDisks {
	return GetDisks{
		vmService: vmService,
	}
}

func (gd GetDisks) Run(vmCID VMCID) (disks []string, err error) {
	disks, err = gd.vmService.AttachedDisks(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Finding disks for vm '%s'", vmCID)
	}

	return disks, nil
}
