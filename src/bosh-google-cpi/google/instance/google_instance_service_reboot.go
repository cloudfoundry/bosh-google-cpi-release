package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
)

const (
	STATUS_RUNNING    = "RUNNING"
	STATUS_TERMINATED = "TERMINATED"
)

func (i GoogleInstanceService) Reboot(id string) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	switch instance.Status {
	default:
		return bosherr.Errorf("Can not reboot instance in state %q", instance.Status)
	case STATUS_RUNNING:
		i.logger.Debug(googleInstanceServiceLogTag, "Rebooting running Google Instance %q via reset API", id)
		operation, err := i.computeService.Instances.Reset(i.project, util.ResourceSplitter(instance.Zone), id).Do()
		if err != nil {
			return bosherr.WrapErrorf(err, "Failed to reboot Google Instance '%s'", id)
		}
		if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
			return bosherr.WrapErrorf(err, "Failed to reboot Google Instance '%s'", id)
		}
		return nil
	case STATUS_TERMINATED:
		i.logger.Debug(googleInstanceServiceLogTag, "Rebooting terminated Google Instance %q via start API", id)
		operation, err := i.computeService.Instances.Start(i.project, util.ResourceSplitter(instance.Zone), id).Do()
		if err != nil {
			return bosherr.WrapErrorf(err, "Failed to reboot Google Instance '%s'", id)
		}
		if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
			return bosherr.WrapErrorf(err, "Failed to reboot Google Instance '%s'", id)
		}
		return nil
	}
}
