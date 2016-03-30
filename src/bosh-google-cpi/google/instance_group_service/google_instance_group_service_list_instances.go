package instancegroup

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceGroupService) ListInstances(id string, zone string) ([]string, error) {
	var instances []string

	instanceGroupsRequest := &compute.InstanceGroupsListInstancesRequest{InstanceState: "ALL"}

	i.logger.Debug(googleInstanceGroupServiceLogTag, "Listing Google Instances for Google Instance Group '%s' in zone '%s'", id, util.ResourceSplitter(zone))
	instancesList, err := i.computeService.InstanceGroups.ListInstances(i.project, util.ResourceSplitter(zone), id, instanceGroupsRequest).Do()
	if err != nil {
		return instances, bosherr.WrapErrorf(err, "Failed to list Google Instances for Google Instance Group '%s' in zone '%s'", id, util.ResourceSplitter(zone))
	}

	for _, instanceItem := range instancesList.Items {
		instances = append(instances, instanceItem.Instance)
	}

	return instances, nil
}
