package acceleratortype

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"
)

const googleAcceleratorTypeServiceLogTag = "GoogleAcceleratorTypeService"

type GoogleAcceleratorTypeService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleAcceleratorTypeService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleAcceleratorTypeService {
	return GoogleAcceleratorTypeService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
