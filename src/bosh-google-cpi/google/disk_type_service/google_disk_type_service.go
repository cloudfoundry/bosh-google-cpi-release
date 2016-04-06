package disktype

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"google.golang.org/api/compute/v1"
)

const googleDiskTypeServiceLogTag = "GoogleDiskTypeService"

type GoogleDiskTypeService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleDiskTypeService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleDiskTypeService {
	return GoogleDiskTypeService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
