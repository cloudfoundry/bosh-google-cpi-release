package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/google/instance_service"
)

type HasVM struct {
	vmService instance.Service
}

func NewHasVM(
	vmService instance.Service,
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
