package targetpool

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (t GoogleTargetPoolService) AddInstance(id string, vmLink string) error {
	targetPool, found, err := t.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "Google Target Pool '%s' does not exists", id)
	}

	for _, vm := range targetPool.Instances {
		if vm == vmLink {
			t.logger.Debug(googleTargetPoolServiceLogTag, "Google Instance '%s' already attached to Google Target Pool '%s'", util.ResourceSplitter(vmLink), id)
			return nil
		}
	}

	var instances []*compute.InstanceReference
	instance := &compute.InstanceReference{Instance: vmLink}
	instances = append(instances, instance)
	targetPoolsRequest := &compute.TargetPoolsAddInstanceRequest{Instances: instances}

	t.logger.Debug(googleTargetPoolServiceLogTag, "Adding Google Instance '%s' to Google Target Pool '%s'", util.ResourceSplitter(vmLink), id)
	operation, err := t.computeService.TargetPools.AddInstance(t.project, util.ResourceSplitter(targetPool.Region), id, targetPoolsRequest).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to add Google Instance '%s' to Target Pool '%s'", util.ResourceSplitter(vmLink), id)
	}

	if _, err = t.operationService.Waiter(operation, "", targetPool.Region); err != nil {
		return bosherr.WrapErrorf(err, "Failed to add Google Instance '%s' to Target Pool '%s'", util.ResourceSplitter(vmLink), id)
	}

	return nil
}
