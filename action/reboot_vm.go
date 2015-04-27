package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/instance"
)

type RebootVM struct {
	vmService ginstance.GoogleInstanceService
}

func NewRebootVM(
	vmService ginstance.GoogleInstanceService,
) RebootVM {
	return RebootVM{
		vmService: vmService,
	}
}

func (rv RebootVM) Run(vmCID VMCID) (interface{}, error) {
	err := rv.vmService.Reboot(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Rebooting vm '%s'", vmCID)
	}

	return nil, nil
}
