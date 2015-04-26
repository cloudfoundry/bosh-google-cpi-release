package ginstance

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
)

func (i GoogleInstanceService) Reboot(id string) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Instance '%s' not found", id)
	}

	i.logger.Debug(googleInstanceServiceLogTag, "Rebooting Google Instance '%s'", id)
	operation, err := i.computeService.Instances.Reset(i.project, gutil.ResourceSplitter(instance.Zone), id).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to reboot Google Instance '%s'", id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to reboot Google Instance '%s'", id)
	}

	return nil
}
