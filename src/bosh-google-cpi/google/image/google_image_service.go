package image

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/storage/v1"

	"bosh-google-cpi/google/operation"
)

const googleImageServiceLogTag = "GoogleImageService"
const googleImageNamePrefix = "stemcell"
const googleImageDescription = "Image managed by BOSH"
const googleImageReadyStatus = "READY"
const googleImageFailedStatus = "FAILED"

type GoogleImageService struct {
	project          string
	computeService   *compute.Service
	storageService   *storage.Service
	operationService operation.Service
	uuidGen          boshuuid.Generator
	logger           boshlog.Logger
}

func NewGoogleImageService(
	project string,
	computeService *compute.Service,
	storageService *storage.Service,
	operationService operation.Service,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) GoogleImageService {
	return GoogleImageService{
		project:          project,
		computeService:   computeService,
		storageService:   storageService,
		operationService: operationService,
		uuidGen:          uuidGen,
		logger:           logger,
	}
}
