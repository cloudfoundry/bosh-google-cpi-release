package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
)

func (i GoogleInstanceService) Delete(id string) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	i.logger.Debug(googleInstanceServiceLogTag, "Deleting Google Instance '%s'", id)
	operation, err := i.computeService.Instances.Delete(i.project, util.ResourceSplitter(instance.Zone), id).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Instance '%s'", id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Instance '%s'", id)
	}

	if err = i.removeFromTargetPool(instance.SelfLink); err != nil {
		return bosherr.WrapErrorf(err, "Failed to remove Google Instance %q from Target Pool", id)
	}

	if err = i.removeFromBackendService(instance.SelfLink); err != nil {
		return bosherr.WrapErrorf(err, "Failed to remove Google Instance %q from Backend Services", id)
	}
	return nil
}
