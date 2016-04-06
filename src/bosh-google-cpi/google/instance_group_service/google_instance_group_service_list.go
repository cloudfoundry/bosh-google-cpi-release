package instancegroup

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
)

func (i GoogleInstanceGroupService) List(zone string) ([]InstanceGroup, error) {
	var instanceGroups []InstanceGroup

	if zone == "" {
		i.logger.Debug(googleInstanceGroupServiceLogTag, "Listing Google Instance Groups")
		instanceGroupAggregatedList, err := i.computeService.InstanceGroups.AggregatedList(i.project).Do()
		if err != nil {
			return instanceGroups, bosherr.WrapError(err, "Failed to list Google Instance Groups")
		}

		for _, instanceGroupList := range instanceGroupAggregatedList.Items {
			for _, instanceGroupItem := range instanceGroupList.InstanceGroups {
				instances, err := i.ListInstances(instanceGroupItem.Name, instanceGroupItem.Zone)
				if err != nil {
					return instanceGroups, bosherr.WrapErrorf(err, "Failed to list Google Instances for Google Instance Group '%s' in zone '%s'", instanceGroupItem.Name, instanceGroupItem.Zone)
				}

				instanceGroup := InstanceGroup{
					Name:      instanceGroupItem.Name,
					Instances: instances,
					SelfLink:  instanceGroupItem.SelfLink,
					Zone:      instanceGroupItem.Zone,
				}
				instanceGroups = append(instanceGroups, instanceGroup)
			}
		}

		return instanceGroups, nil
	}

	i.logger.Debug(googleInstanceGroupServiceLogTag, "Listing Google Instance Groups in zone '%s'", zone)
	instanceGroupList, err := i.computeService.InstanceGroups.List(i.project, util.ResourceSplitter(zone)).Do()
	if err != nil {
		return instanceGroups, bosherr.WrapErrorf(err, "Failed to list Google Instance Groups in zone '%s'", zone)
	}

	for _, instanceGroupItem := range instanceGroupList.Items {
		instances, err := i.ListInstances(instanceGroupItem.Name, instanceGroupItem.Zone)
		if err != nil {
			return instanceGroups, bosherr.WrapErrorf(err, "Failed to list Google Instances for Google Instance Group '%s' in zone '%s'", instanceGroupItem.Name, instanceGroupItem.Zone)
		}

		instanceGroup := InstanceGroup{
			Name:      instanceGroupItem.Name,
			Instances: instances,
			SelfLink:  instanceGroupItem.SelfLink,
			Zone:      instanceGroupItem.Zone,
		}
		instanceGroups = append(instanceGroups, instanceGroup)
	}

	return instanceGroups, nil
}
