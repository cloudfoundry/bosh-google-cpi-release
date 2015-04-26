package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/image"
)

type DeleteStemcell struct {
	stemcellService gimage.GoogleImageService
}

func NewDeleteStemcell(
	stemcellService gimage.GoogleImageService,
) DeleteStemcell {
	return DeleteStemcell{
		stemcellService: stemcellService,
	}
}

func (ds DeleteStemcell) Run(stemcellCID StemcellCID) (interface{}, error) {
	err := ds.stemcellService.Delete(string(stemcellCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting stemcell '%s'", stemcellCID)
	}

	return nil, nil
}
