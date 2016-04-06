package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/google/image_service"
)

type DeleteStemcell struct {
	imageService image.Service
}

func NewDeleteStemcell(
	imageService image.Service,
) DeleteStemcell {
	return DeleteStemcell{
		imageService: imageService,
	}
}

func (ds DeleteStemcell) Run(stemcellCID StemcellCID) (interface{}, error) {
	if err := ds.imageService.Delete(string(stemcellCID)); err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting stemcell '%s'", stemcellCID)
	}

	return nil, nil
}
