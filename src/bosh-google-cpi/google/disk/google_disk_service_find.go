package disk

import (
	"context"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"bosh-google-cpi/util"
)

func (d GoogleDiskService) Find(id string, zone string) (Disk, bool, error) {
	if zone == "" {
		d.logger.Debug(googleDiskServiceLogTag, "Finding Google Disk '%s'", id)
		filter := fmt.Sprintf("name = %q", id)
		var found *compute.Disk
		err := d.computeService.Disks.AggregatedList(d.project).Filter(filter).Pages(context.Background(), func(page *compute.DiskAggregatedList) error {
			for _, items := range page.Items {
				for _, item := range items.Disks {
					if item.Name != id {
						continue
					}
					if found != nil {
						return fmt.Errorf("disk %q is ambiguous across zones", id)
					}
					found = item
				}
			}
			return nil
		})
		if err != nil {
			return Disk{}, false, bosherr.WrapErrorf(err, "Failed to find Google Disk '%s'", id)
		}
		if found != nil {
			return Disk{Name: found.Name, SelfLink: found.SelfLink, Status: found.Status, Zone: found.Zone, Users: found.Users, SizeGb: found.SizeGb, Type: found.Type}, true, nil
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
		SizeGb:   diskItem.SizeGb,
		Type:     diskItem.Type,
	}
	return disk, true, nil
}
