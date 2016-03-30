package disktype

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/googleapi"
)

func (d GoogleDiskTypeService) Find(id string, zone string) (DiskType, bool, error) {
	d.logger.Debug(googleDiskTypeServiceLogTag, "Finding Google Disk Type '%s' in zone '%s'", id, zone)
	diskTypeItem, err := d.computeService.DiskTypes.Get(d.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return DiskType{}, false, nil
		}

		return DiskType{}, false, bosherr.WrapErrorf(err, "Failed to find Google Disk Type '%s' in zone '%s'", id, zone)
	}

	diskType := DiskType{
		Name:     diskTypeItem.Name,
		SelfLink: diskTypeItem.SelfLink,
		Zone:     diskTypeItem.Zone,
	}
	return diskType, true, nil
}
