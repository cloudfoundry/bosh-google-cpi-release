package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/disk_service"
)

type DeleteDisk struct {
	diskService disk.Service
}

func NewDeleteDisk(
	diskService disk.Service,
) DeleteDisk {
	return DeleteDisk{
		diskService: diskService,
	}
}

func (dd DeleteDisk) Run(diskCID DiskCID) (interface{}, error) {
	err := dd.diskService.Delete(string(diskCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting disk '%s'", diskCID)
	}

	return nil, nil
}
