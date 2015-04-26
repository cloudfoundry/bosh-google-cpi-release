package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/instance"
)

type SetVMMetadata struct {
	vmService ginstance.GoogleInstanceService
}

func NewSetVMMetadata(
	vmService ginstance.GoogleInstanceService,
) SetVMMetadata {
	return SetVMMetadata{
		vmService: vmService,
	}
}

func (svm SetVMMetadata) Run(vmCID VMCID, vmMetadata VMMetadata) (interface{}, error) {
	err := svm.vmService.SetMetadata(string(vmCID), ginstance.GoogleInstanceMetadata(vmMetadata))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Setting metadata for vm '%s'", vmCID)
	}

	return nil, nil
}
