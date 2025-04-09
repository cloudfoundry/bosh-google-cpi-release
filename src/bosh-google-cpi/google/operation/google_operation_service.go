package operation

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

const googleOperationServiceLogTag = "GoogleOperationService"
const googleOperationServiceMaxTries = 100
const googleOperationServiceMaxSleepExponent = 3
const googleOperationReadyStatus = "DONE"

type GoogleOperationService struct {
	project         string
	computeService  *compute.Service
	computeServiceB *computebeta.Service
	logger          boshlog.Logger
}

func NewGoogleOperationService(
	project string,
	computeService *compute.Service,
	computeServiceB *computebeta.Service,
	logger boshlog.Logger,
) GoogleOperationService {
	return GoogleOperationService{
		project:         project,
		computeService:  computeService,
		computeServiceB: computeServiceB,
		logger:          logger,
	}
}
