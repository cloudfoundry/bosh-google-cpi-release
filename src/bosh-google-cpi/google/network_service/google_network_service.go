package network

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"

	project "bosh-google-cpi/google/project_service"
)

const googleNetworkServiceLogTag = "GoogleNetworkService"

type GoogleNetworkService struct {
	projectService project.Service
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleNetworkService(
	projectService project.Service,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleNetworkService {
	return GoogleNetworkService{
		projectService: projectService,
		computeService: computeService,
		logger:         logger,
	}
}
