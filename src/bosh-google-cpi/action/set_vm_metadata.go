package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/instance_service"
)

type SetVMMetadata struct {
	vmService instance.Service
}

func NewSetVMMetadata(
	vmService instance.Service,
) SetVMMetadata {
	return SetVMMetadata{
		vmService: vmService,
	}
}

func (svm SetVMMetadata) Run(vmCID VMCID, vmMetadata VMMetadata) (interface{}, error) {
	if err := svm.vmService.SetMetadata(string(vmCID), instance.Metadata(vmMetadata)); err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Setting metadata for vm '%s'", vmCID)
	}

	return nil, nil
}
