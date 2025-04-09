package subnetwork

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"

	"bosh-google-cpi/google/project"
)

const googleSubnetworkServiceLogTag = "GoogleSubnetworkService"

type GoogleSubnetworkService struct {
	projectService project.Service
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleSubnetworkService(
	projectService project.Service,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleSubnetworkService {
	return GoogleSubnetworkService{
		projectService: projectService,
		computeService: computeService,
		logger:         logger,
	}
}
