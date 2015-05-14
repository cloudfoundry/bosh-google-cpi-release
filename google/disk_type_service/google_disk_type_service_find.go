package gdisktype

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (d GoogleDiskTypeService) Find(id string, zone string) (*compute.DiskType, bool, error) {
	d.logger.Debug(googleDiskTypeServiceLogTag, "Finding Google Disk Type '%s' in zone '%s'", id, zone)
	diskType, err := d.computeService.DiskTypes.Get(d.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.DiskType{}, false, nil
		}

		return &compute.DiskType{}, false, bosherr.WrapErrorf(err, "Failed to find Google Disk Type '%s' in zone '%s'", id, zone)
	}

	return diskType, true, nil
}
