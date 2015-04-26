package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/instance"
)

type HasVM struct {
	vmService ginstance.GoogleInstanceService
}

func NewHasVM(
	vmService ginstance.GoogleInstanceService,
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
