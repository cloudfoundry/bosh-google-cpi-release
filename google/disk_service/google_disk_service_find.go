package gdisk

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (d GoogleDiskService) Find(id string, zone string) (*compute.Disk, bool, error) {
	if zone == "" {
		d.logger.Debug(googleDiskServiceLogTag, "Finding Google Disk '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		disks, err := d.computeService.Disks.AggregatedList(d.project).Filter(filter).Do()
		if err != nil {
			return &compute.Disk{}, false, bosherr.WrapErrorf(err, "Failed to find Google Disk '%s'", id)
		}

		for _, diskItems := range disks.Items {
			for _, disk := range diskItems.Disks {
				// Return the first disk (it can only be 1 disk with the same name across all zones)
				return disk, true, nil
			}
		}

		return &compute.Disk{}, false, nil
	}

	d.logger.Debug(googleDiskServiceLogTag, "Finding Google Disk '%s' in zone '%s'", id, zone)
	disk, err := d.computeService.Disks.Get(d.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.Disk{}, false, nil
		}

		return &compute.Disk{}, false, bosherr.WrapErrorf(err, "Failed to find Google Disk '%s' in zone '%s'", id, util.ResourceSplitter(zone))
	}

	return disk, true, nil
}
