package instancegroup

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"bosh-google-cpi/google/operation_service"
	"google.golang.org/api/compute/v1"
)

const googleInstanceGroupServiceLogTag = "GoogleInstanceGroupService"

type GoogleInstanceGroupService struct {
	project          string
	computeService   *compute.Service
	operationService operation.Service
	logger           boshlog.Logger
}

func NewGoogleInstanceGroupService(
	project string,
	computeService *compute.Service,
	operationService operation.Service,
	logger boshlog.Logger,
) GoogleInstanceGroupService {
	return GoogleInstanceGroupService{
		project:          project,
		computeService:   computeService,
		operationService: operationService,
		logger:           logger,
	}
}
