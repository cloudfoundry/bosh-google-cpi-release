package subnetwork

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"bosh-google-cpi/google/project_service"

	"google.golang.org/api/compute/v1"
)

const googleSubnetworkServiceLogTag = "GoogleSubnetworkService"

type GoogleSubnetworkService struct {
	projcetService project.Service
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleSubnetworkService(
	projcetService project.Service,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleSubnetworkService {
	return GoogleSubnetworkService{
		projcetService: projcetService,
		computeService: computeService,
		logger:         logger,
	}
}
