package action

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/google/image"
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
	if strings.HasPrefix(string(stemcellCID), "https://www.googleapis.com/compute/v1/projects/") {
		return nil, nil
	}

	if err := ds.imageService.Delete(string(stemcellCID)); err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting stemcell '%s'", stemcellCID)
	}

	return nil, nil
}
