package gimage

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (i GoogleImageService) Find(id string) (*compute.Image, bool, error) {
	i.logger.Debug(googleImageServiceLogTag, "Finding Google Image '%s'", id)
	image, err := i.computeService.Images.Get(i.project, id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.Image{}, false, nil
		}

		return &compute.Image{}, false, bosherr.WrapErrorf(err, "Failed to find Google Image '%s'", id)
	}

	return image, true, nil
}
