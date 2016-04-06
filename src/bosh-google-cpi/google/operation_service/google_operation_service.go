package operation

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"google.golang.org/api/compute/v1"
)

const googleOperationServiceLogTag = "GoogleOperationService"
const googleOperationServiceMaxTries = 100
const googleOperationServiceMaxSleepExponent = 3
const googleOperationReadyStatus = "DONE"

type GoogleOperationService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleOperationService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleOperationService {
	return GoogleOperationService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
