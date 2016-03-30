package instance

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"bosh-google-cpi/google/address_service"
	"bosh-google-cpi/google/instance_group_service"
	"bosh-google-cpi/google/network_service"
	"bosh-google-cpi/google/operation_service"
	"bosh-google-cpi/google/subnetwork_service"
	"bosh-google-cpi/google/target_pool_service"
	"google.golang.org/api/compute/v1"
)

const googleInstanceServiceLogTag = "GoogleInstanceService"
const googleInstanceNamePrefix = "vm"
const googleInstanceDescription = "Instance managed by BOSH"

type GoogleInstanceService struct {
	project              string
	computeService       *compute.Service
	addressService       address.Service
	instanceGroupService instancegroup.Service
	networkService       network.Service
	operationService     operation.Service
	subnetworkService    subnetwork.Service
	targetPoolService    targetpool.Service
	uuidGen              boshuuid.Generator
	logger               boshlog.Logger
}

func NewGoogleInstanceService(
	project string,
	computeService *compute.Service,
	addressService address.Service,
	instanceGroupService instancegroup.Service,
	networkService network.Service,
	operationService operation.Service,
	subnetworkService subnetwork.Service,
	targetPoolService targetpool.Service,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) GoogleInstanceService {
	return GoogleInstanceService{
		project:              project,
		computeService:       computeService,
		addressService:       addressService,
		instanceGroupService: instanceGroupService,
		networkService:       networkService,
		operationService:     operationService,
		subnetworkService:    subnetworkService,
		targetPoolService:    targetPoolService,
		uuidGen:              uuidGen,
		logger:               logger,
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
