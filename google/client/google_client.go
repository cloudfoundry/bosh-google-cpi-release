package gclient

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"golang.org/x/oauth2"
	oauthgoogle "golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/storage/v1"
)

const computeScope = compute.ComputeScope
const storageScope = storage.DevstorageFull_controlScope

type GoogleClient struct {
	project        string
	jsonKey        string
	defaultZone    string
	computeService *compute.Service
	storageService *storage.Service
}

func NewGoogleClient(
	project string,
	jsonKey string,
	defaultZone string,
) (GoogleClient, error) {
	computeJwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(jsonKey), computeScope)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Reading Google JSON Key")
	}

	computeClient := computeJwtConf.Client(oauth2.NoContext)
	computeService, err := compute.New(computeClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Compute Service client")
	}

	storageJwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(jsonKey), storageScope)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Reading Google JSON Key")
	}

	storageClient := storageJwtConf.Client(oauth2.NoContext)
	storageService, err := storage.New(storageClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Storage Service client")
	}

	return GoogleClient{
		project:        project,
		jsonKey:        jsonKey,
		defaultZone:    defaultZone,
		computeService: computeService,
		storageService: storageService,
	}, nil
}

func (c GoogleClient) Project() string {
	return c.project
}

func (c GoogleClient) DefaultZone() string {
	return c.defaultZone
}

func (c GoogleClient) ComputeService() *compute.Service {
	return c.computeService
}

func (c GoogleClient) StorageService() *storage.Service {
	return c.storageService
}
