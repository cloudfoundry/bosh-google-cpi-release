package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceService) AddAccessConfig(id string, zone string, networkInterface string, accessConfig *compute.AccessConfig) error {
	i.logger.Debug(googleInstanceServiceLogTag, "Adding access config for Google Instance '%s'", id)
	operation, err := i.computeService.Instances.AddAccessConfig(i.project, util.ResourceSplitter(zone), id, networkInterface, accessConfig).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to add access config for Google Instance '%s'", id)
	}

	if _, err = i.operationService.Waiter(operation, zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to add access config for Google Instance '%s'", id)
	}

	return nil
}
