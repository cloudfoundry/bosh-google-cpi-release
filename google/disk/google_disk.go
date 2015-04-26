package gdisk

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	"github.com/frodenas/bosh-google-cpi/google/operation"
	"google.golang.org/api/compute/v1"
)

const googleDiskServiceLogTag = "GoogleDiskService"
const googleDiskNamePrefix = "disk"
const googleDiskDescription = "Disk managed by BOSH"
const googleDiskReadyStatus = "READY"

type GoogleDiskService struct {
	project          string
	computeService   *compute.Service
	operationService goperation.GoogleOperationService
	uuidGen          boshuuid.Generator
	logger           boshlog.Logger
}

func NewGoogleDiskService(
	project string,
	computeService *compute.Service,
	operationService goperation.GoogleOperationService,
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
