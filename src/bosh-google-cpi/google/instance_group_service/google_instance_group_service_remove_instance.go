package instancegroup

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceGroupService) RemoveInstance(id string, vmLink string) error {
	instanceGroup, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "Google Instance Group '%s' does not exists", id)
	}

	attached := false
	for _, vm := range instanceGroup.Instances {
		if vm == vmLink {
			attached = true
		}
	}
	if !attached {
		i.logger.Debug(googleInstanceGroupServiceLogTag, "Google Instance '%s' is not attached to Google Instance Group '%s'", util.ResourceSplitter(vmLink), id)
		return nil
	}

	var instances []*compute.InstanceReference
	instance := &compute.InstanceReference{Instance: vmLink}
	instances = append(instances, instance)
	instanceGroupsRequest := &compute.InstanceGroupsRemoveInstancesRequest{Instances: instances}

	i.logger.Debug(googleInstanceGroupServiceLogTag, "Removing Google Instance '%s' from Google Instance Group '%s'", util.ResourceSplitter(vmLink), id)
	operation, err := i.computeService.InstanceGroups.RemoveInstances(i.project, util.ResourceSplitter(instanceGroup.Zone), id, instanceGroupsRequest).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to remove Google Instance '%s' from Google Instance Group '%s'", util.ResourceSplitter(vmLink), id)
	}

	if _, err = i.operationService.Waiter(operation, instanceGroup.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to remove Google Instance '%s' from Google Instance Group '%s'", util.ResourceSplitter(vmLink), id)
	}

	return nil
}
