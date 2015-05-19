package image

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"google.golang.org/api/googleapi"
)

func (i GoogleImageService) Find(id string) (Image, bool, error) {
	i.logger.Debug(googleImageServiceLogTag, "Finding Google Image '%s'", id)
	imageItem, err := i.computeService.Images.Get(i.project, id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return Image{}, false, nil
		}

		return Image{}, false, bosherr.WrapErrorf(err, "Failed to find Google Image '%s'", id)
	}

	image := Image{
		Name:     imageItem.Name,
		SelfLink: imageItem.SelfLink,
		Status:   imageItem.Status,
	}
	return image, true, nil
}
