package targetpool

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (t GoogleTargetPoolService) RemoveInstance(id string, vmLink string) error {
	targetPool, found, err := t.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "Google Target Pool '%s' does not exists", id)
	}

	attached := false
	for _, vm := range targetPool.Instances {
		if vm == vmLink {
			attached = true
		}
	}
	if !attached {
		t.logger.Debug(googleTargetPoolServiceLogTag, "Google Instance '%s' is not attached to Google Target Pool '%s'", util.ResourceSplitter(vmLink), id)
		return nil
	}

	var instances []*compute.InstanceReference
	instance := &compute.InstanceReference{Instance: vmLink}
	instances = append(instances, instance)
	targetPoolsRequest := &compute.TargetPoolsRemoveInstanceRequest{Instances: instances}

	t.logger.Debug(googleTargetPoolServiceLogTag, "Removing Google Instance '%s' from Google Target Pool '%s'", util.ResourceSplitter(vmLink), id)
	operation, err := t.computeService.TargetPools.RemoveInstance(t.project, util.ResourceSplitter(targetPool.Region), id, targetPoolsRequest).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to remove Google Instance '%s' from Target Pool '%s'", util.ResourceSplitter(vmLink), id)
	}

	if _, err = t.operationService.Waiter(operation, "", targetPool.Region); err != nil {
		return bosherr.WrapErrorf(err, "Failed to remove Google Instance '%s' from Target Pool '%s'", util.ResourceSplitter(vmLink), id)
	}

	return nil
}
