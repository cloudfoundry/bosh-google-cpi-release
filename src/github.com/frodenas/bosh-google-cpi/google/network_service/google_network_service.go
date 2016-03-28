package network

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"google.golang.org/api/compute/v1"
)

const googleNetworkServiceLogTag = "GoogleNetworkService"

type GoogleNetworkService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleNetworkService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleNetworkService {
	return GoogleNetworkService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
