package instancegroup

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/googleapi"
)

func (i GoogleInstanceGroupService) Find(id string, zone string) (InstanceGroup, bool, error) {
	if zone == "" {
		i.logger.Debug(googleInstanceGroupServiceLogTag, "Finding Google Instance Group '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		instanceGroups, err := i.computeService.InstanceGroups.AggregatedList(i.project).Filter(filter).Do()
		if err != nil {
			return InstanceGroup{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance Group '%s'", id)
		}

		for _, instanceGroupItems := range instanceGroups.Items {
			for _, instanceGroupItem := range instanceGroupItems.InstanceGroups {
				instances, err := i.ListInstances(id, instanceGroupItem.Zone)
				if err != nil {
					return InstanceGroup{}, false, bosherr.WrapErrorf(err, "Failed to list Google Instances for Google Instance Group '%s' in zone '%s", id, zone)
				}

				// Return the first instance group (it can only be 1 instance group with the same name across all zones)
				instanceGroup := InstanceGroup{
					Name:      instanceGroupItem.Name,
					Instances: instances,
					SelfLink:  instanceGroupItem.SelfLink,
					Zone:      instanceGroupItem.Zone,
				}
				return instanceGroup, true, nil
			}
		}

		return InstanceGroup{}, false, nil
	}

	i.logger.Debug(googleInstanceGroupServiceLogTag, "Finding Google Instance Group '%s' in zone '%s'", id, zone)
	instanceGroupItem, err := i.computeService.InstanceGroups.Get(i.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return InstanceGroup{}, false, nil
		}

		return InstanceGroup{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance Group '%s' in zone '%s'", id, zone)
	}

	instances, err := i.ListInstances(id, zone)
	if err != nil {
		return InstanceGroup{}, false, bosherr.WrapErrorf(err, "Failed to list Google Instances for Google Instance Group '%s' in zone '%s'", id, zone)
	}

	instanceGroup := InstanceGroup{
		Name:      instanceGroupItem.Name,
		Instances: instances,
		SelfLink:  instanceGroupItem.SelfLink,
		Zone:      instanceGroupItem.Zone,
	}
	return instanceGroup, true, nil
}
