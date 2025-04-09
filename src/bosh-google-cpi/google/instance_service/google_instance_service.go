package instance

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"

	address "bosh-google-cpi/google/address_service"
	backendservice "bosh-google-cpi/google/backendservice_service"
	disktype "bosh-google-cpi/google/disk_type_service"
	network "bosh-google-cpi/google/network_service"
	operation "bosh-google-cpi/google/operation_service"
	subnetwork "bosh-google-cpi/google/subnetwork_service"
	targetpool "bosh-google-cpi/google/target_pool_service"
)

const googleInstanceServiceLogTag = "GoogleInstanceService"
const googleInstanceNamePrefix = "vm"
const googleInstanceDescription = "Instance managed by BOSH"

type GoogleInstanceService struct {
	project               string
	computeService        *compute.Service
	computeServiceB       *computebeta.Service
	addressService        address.Service
	backendServiceService backendservice.Service
	networkService        network.Service
	operationService      operation.Service
	diskTypeService       disktype.Service
	subnetworkService     subnetwork.Service
	targetPoolService     targetpool.Service
	uuidGen               boshuuid.Generator
	logger                boshlog.Logger
}

func NewGoogleInstanceService(
	project string,
	computeService *compute.Service,
	computeServiceB *computebeta.Service,
	addressService address.Service,
	backendServiceService backendservice.Service,
	networkService network.Service,
	operationService operation.Service,
	subnetworkService subnetwork.Service,
	targetPoolService targetpool.Service,
	diskTypeService disktype.Service,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) GoogleInstanceService {
	return GoogleInstanceService{
		project:               project,
		computeService:        computeService,
		computeServiceB:       computeServiceB,
		addressService:        addressService,
		backendServiceService: backendServiceService,
		networkService:        networkService,
		operationService:      operationService,
		subnetworkService:     subnetworkService,
		targetPoolService:     targetPoolService,
		diskTypeService:       diskTypeService,
		uuidGen:               uuidGen,
		logger:                logger,
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
