package machinetype

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"google.golang.org/api/compute/v1"
)

const googleMachineTypeServiceLogTag = "GoogleMachineTypeService"

type GoogleMachineTypeService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleMachineTypeService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleMachineTypeService {
	return GoogleMachineTypeService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
