package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-google-cpi/google/image_service"
)

type DeleteStemcell struct {
	stemcellService image.Service
}

func NewDeleteStemcell(
	stemcellService image.Service,
) DeleteStemcell {
	return DeleteStemcell{
		stemcellService: stemcellService,
	}
}

func (ds DeleteStemcell) Run(stemcellCID StemcellCID) (interface{}, error) {
	if err := ds.stemcellService.Delete(string(stemcellCID)); err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting stemcell '%s'", stemcellCID)
	}

	return nil, nil
}
