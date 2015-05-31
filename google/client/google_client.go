package client

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-google-cpi/google/config"

	"golang.org/x/oauth2"
	oauthgoogle "golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/storage/v1"
)

const computeScope = compute.ComputeScope
const storageScope = storage.DevstorageFullControlScope

type GoogleClient struct {
	config         config.Config
	computeService *compute.Service
	storageService *storage.Service
}

func NewGoogleClient(
	config config.Config,
) (GoogleClient, error) {
	computeJwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(config.JSONKey), computeScope)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Reading Google JSON Key")
	}

	computeClient := computeJwtConf.Client(oauth2.NoContext)
	computeService, err := compute.New(computeClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Compute Service client")
	}

	storageJwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(config.JSONKey), storageScope)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Reading Google JSON Key")
	}

	storageClient := storageJwtConf.Client(oauth2.NoContext)
	storageService, err := storage.New(storageClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Storage Service client")
	}

	return GoogleClient{
		config:         config,
		computeService: computeService,
		storageService: storageService,
	}, nil
}

func (c GoogleClient) Project() string {
	return c.config.Project
}

func (c GoogleClient) DefaultRootDiskSizeGb() int {
	return c.config.DefaultRootDiskSizeGb
}

func (c GoogleClient) DefaultRootDiskType() string {
	return c.config.DefaultRootDiskType
}

func (c GoogleClient) DefaultZone() string {
	return c.config.DefaultZone
}

func (c GoogleClient) ComputeService() *compute.Service {
	return c.computeService
}

func (c GoogleClient) StorageService() *storage.Service {
	return c.storageService
}
