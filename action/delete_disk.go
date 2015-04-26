package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/disk"
)

type DeleteDisk struct {
	diskService gdisk.GoogleDiskService
}

func NewDeleteDisk(
	diskService gdisk.GoogleDiskService,
) DeleteDisk {
	return DeleteDisk{
		diskService: diskService,
	}
}

func (dd DeleteDisk) Run(diskCID DiskCID) (interface{}, error) {
	err := dd.diskService.Delete(string(diskCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting disk '%s'", diskCID)
	}

	return nil, nil
}
