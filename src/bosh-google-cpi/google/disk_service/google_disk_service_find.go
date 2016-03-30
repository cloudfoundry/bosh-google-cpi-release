package disk

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/googleapi"
)

func (d GoogleDiskService) Find(id string, zone string) (Disk, bool, error) {
	if zone == "" {
		d.logger.Debug(googleDiskServiceLogTag, "Finding Google Disk '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		disks, err := d.computeService.Disks.AggregatedList(d.project).Filter(filter).Do()
		if err != nil {
			return Disk{}, false, bosherr.WrapErrorf(err, "Failed to find Google Disk '%s'", id)
		}

		for _, diskItems := range disks.Items {
			for _, diskItem := range diskItems.Disks {
				// Return the first disk (it can only be 1 disk with the same name across all zones)
				disk := Disk{
					Name:     diskItem.Name,
					SelfLink: diskItem.SelfLink,
					Status:   diskItem.Status,
					Zone:     diskItem.Zone,
				}
				return disk, true, nil
			}
		}

		return Disk{}, false, nil
	}

	d.logger.Debug(googleDiskServiceLogTag, "Finding Google Disk '%s' in zone '%s'", id, zone)
	diskItem, err := d.computeService.Disks.Get(d.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return Disk{}, false, nil
		}

		return Disk{}, false, bosherr.WrapErrorf(err, "Failed to find Google Disk '%s' in zone '%s'", id, util.ResourceSplitter(zone))
	}

	disk := Disk{
		Name:     diskItem.Name,
		SelfLink: diskItem.SelfLink,
		Status:   diskItem.Status,
		Zone:     diskItem.Zone,
	}
	return disk, true, nil
}
