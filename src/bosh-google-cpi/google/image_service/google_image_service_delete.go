package image

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (i GoogleImageService) Delete(id string) error {
	image, found, err := i.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "Google Image '%s' does not exists", id)
	}

	if image.Status != googleImageReadyStatus && image.Status != googleImageFailedStatus {
		return bosherr.WrapErrorf(err, "Cannot delete Google Image '%s', status is '%s'", id, image.Status)
	}

	i.logger.Debug(googleImageServiceLogTag, "Deleting Google Image '%s'", id)
	operation, err := i.computeService.Images.Delete(i.project, id).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Image '%s'", id)
	}

	if _, err = i.operationService.Waiter(operation, "", ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Google Image '%s'", id)
	}

	return nil
}
