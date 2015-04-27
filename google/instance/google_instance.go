package ginstance

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	"github.com/frodenas/bosh-google-cpi/google/operation"
	"google.golang.org/api/compute/v1"
)

const googleInstanceServiceLogTag = "GoogleInstanceService"
const googleInstanceNamePrefix = "vm"
const googleInstanceDescription = "Instance managed by BOSH"

type GoogleInstanceService struct {
	project          string
	computeService   *compute.Service
	operationService goperation.GoogleOperationService
	uuidGen          boshuuid.Generator
	logger           boshlog.Logger
}

func NewGoogleInstanceService(
	project string,
	computeService *compute.Service,
	operationService goperation.GoogleOperationService,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) GoogleInstanceService {
	return GoogleInstanceService{
		project:          project,
		computeService:   computeService,
		operationService: operationService,
		uuidGen:          uuidGen,
		logger:           logger,
	}
}

type GoogleInstanceAttachedDisks []string

type GoogleInstanceMetadata map[string]interface{}

type GoogleInstanceProperties struct {
	Zone              string
	Stemcell          string
	MachineType       string
	AutomaticRestart  bool
	OnHostMaintenance string
	ServiceScopes     GoogleInstanceServiceScopes
}

type GoogleInstanceServiceScopes []string

type GoogleUserData struct {
	Instance GoogleUserDataInstanceName     `json:"instance"`
	Registry GoogleUserDataRegistryEndpoint `json:"registry"`
	DNS      GoogleUserDataDNSItems         `json:"dns,omitempty"`
}

type GoogleUserDataInstanceName struct {
	Name string `json:"name"`
}

type GoogleUserDataRegistryEndpoint struct {
	Endpoint string `json:"endpoint"`
}

type GoogleUserDataDNSItems struct {
	NameServers []string `json:"nameserver,omitempty"`
}
