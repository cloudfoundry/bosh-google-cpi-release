package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceService) SetTags(id string, zone string, instanceTags *compute.Tags) error {
	i.logger.Debug(googleInstanceServiceLogTag, "Setting tags for Google Instance '%s'", id)
	operation, err := i.computeService.Instances.SetTags(i.project, util.ResourceSplitter(zone), id, instanceTags).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to set tags for Google Instance '%s'", id)
	}

	if _, err = i.operationService.Waiter(operation, zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to set tags for Google Instance '%s'", id)
	}

	return nil
}
