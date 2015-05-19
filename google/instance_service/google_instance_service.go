package instance

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/frodenas/bosh-google-cpi/google/address_service"
	"github.com/frodenas/bosh-google-cpi/google/network_service"
	"github.com/frodenas/bosh-google-cpi/google/operation_service"
	"github.com/frodenas/bosh-google-cpi/google/target_pool_service"
	"google.golang.org/api/compute/v1"
)

const googleInstanceServiceLogTag = "GoogleInstanceService"
const googleInstanceNamePrefix = "vm"
const googleInstanceDescription = "Instance managed by BOSH"

type GoogleInstanceService struct {
	project           string
	computeService    *compute.Service
	addressService    address.Service
	networkService    network.Service
	operationService  operation.Service
	targetPoolService targetpool.Service
	uuidGen           boshuuid.Generator
	logger            boshlog.Logger
}

func NewGoogleInstanceService(
	project string,
	computeService *compute.Service,
	addressService address.Service,
	networkService network.Service,
	operationService operation.Service,
	targetPoolService targetpool.Service,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) GoogleInstanceService {
	return GoogleInstanceService{
		project:           project,
		computeService:    computeService,
		addressService:    addressService,
		networkService:    networkService,
		operationService:  operationService,
		targetPoolService: targetPoolService,
		uuidGen:           uuidGen,
		logger:            logger,
	}
}

type GoogleUserData struct {
	Server   GoogleUserDataServerName       `json:"server"`
	Registry GoogleUserDataRegistryEndpoint `json:"registry"`
	DNS      GoogleUserDataDNSItems         `json:"dns,omitempty"`
}

type GoogleUserDataServerName struct {
	Name string `json:"name"`
}

type GoogleUserDataRegistryEndpoint struct {
	Endpoint string `json:"endpoint"`
}

type GoogleUserDataDNSItems struct {
	NameServer []string `json:"nameserver,omitempty"`
}
