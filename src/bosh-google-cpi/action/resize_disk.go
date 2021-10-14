package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	disk "bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/util"
)

type ResizeDisk struct {
	diskService disk.Service
}

func NewResizeDisk(
	diskService disk.Service,
) ResizeDisk {
	return ResizeDisk{
		diskService: diskService,
	}
}

func (rd ResizeDisk) Run(diskCID DiskCID, newSize int) (interface{}, error) {
	if err := rd.diskService.Resize(string(diskCID), util.ConvertMib2Gib(newSize)); err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Resizing disk '%s' to '%#v'", diskCID, newSize)
	}

	return nil, nil
}
