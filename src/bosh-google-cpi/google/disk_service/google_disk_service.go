package disk

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"bosh-google-cpi/google/operation_service"
	"google.golang.org/api/compute/v1"
)

const googleDiskServiceLogTag = "GoogleDiskService"
const googleDiskNamePrefix = "disk"
const googleDiskDescription = "Disk managed by BOSH"
const googleDiskReadyStatus = "READY"
const googleDiskFailedStatus = "FAILED"

type GoogleDiskService struct {
	project          string
	computeService   *compute.Service
	operationService operation.Service
	uuidGen          boshuuid.Generator
	logger           boshlog.Logger
}

func NewGoogleDiskService(
	project string,
	computeService *compute.Service,
	operationService operation.Service,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) GoogleDiskService {
	return GoogleDiskService{
		project:          project,
		computeService:   computeService,
		operationService: operationService,
		uuidGen:          uuidGen,
		logger:           logger,
	}
}
