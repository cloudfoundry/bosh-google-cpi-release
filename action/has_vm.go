package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/instance_service"
)

type HasVM struct {
	vmService ginstance.InstanceService
}

func NewHasVM(
	vmService ginstance.InstanceService,
) HasVM {
	return HasVM{
		vmService: vmService,
	}
}

func (hv HasVM) Run(vmCID VMCID) (bool, error) {
	_, found, err := hv.vmService.Find(string(vmCID), "")
	if err != nil {
		return false, bosherr.WrapErrorf(err, "Finding vm '%s'", vmCID)
	}

	return found, nil
}
