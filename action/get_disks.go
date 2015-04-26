package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/instance"
)

type GetDisks struct {
	vmService ginstance.GoogleInstanceService
}

func NewGetDisks(
	vmService ginstance.GoogleInstanceService,
) GetDisks {
	return GetDisks{
		vmService: vmService,
	}
}

func (gd GetDisks) Run(vmCID VMCID) ([]string, error) {
	disks, err := gd.vmService.AttachedDisks(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Finding disks for vm '%s'", vmCID)
	}

	return disks, nil
}
