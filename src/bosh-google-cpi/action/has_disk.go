package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/google/disk_service"
)

type HasDisk struct {
	diskService disk.Service
}

func NewHasDisk(
	diskService disk.Service,
) HasDisk {
	return HasDisk{
		diskService: diskService,
	}
}

func (hd HasDisk) Run(diskCID DiskCID) (bool, error) {
	_, found, err := hd.diskService.Find(string(diskCID), "")
	if err != nil {
		return false, bosherr.WrapErrorf(err, "Finding disk '%s'", diskCID)
	}

	return found, nil
}
