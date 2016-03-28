package subnetwork

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"google.golang.org/api/compute/v1"
)

const googleSubnetworkServiceLogTag = "GoogleSubnetworkService"

type GoogleSubnetworkService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleSubnetworkService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleSubnetworkService {
	return GoogleSubnetworkService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
