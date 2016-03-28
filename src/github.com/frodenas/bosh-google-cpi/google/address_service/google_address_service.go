package address

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"google.golang.org/api/compute/v1"
)

const googleAddressServiceLogTag = "GoogleAddressService"

type GoogleAddressService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleAddressService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleAddressService {
	return GoogleAddressService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
