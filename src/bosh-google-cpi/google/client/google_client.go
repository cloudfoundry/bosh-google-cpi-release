package client

import (
	"net/http"
	"os"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"bosh-google-cpi/google/config"

	"golang.org/x/oauth2"
	oauthgoogle "golang.org/x/oauth2/google"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/storage/v1"
)

const (
	computeScope = compute.ComputeScope
	storageScope = storage.DevstorageFullControlScope
	// Metadata Host needs to be IP address, rather than FQDN, in case the system
	// is set up to use public DNS servers, which would not resolve correctly.
	metadataHost = "169.254.169.254"
)

type GoogleClient struct {
	config          config.Config
	computeService  *compute.Service
	computeServiceB *computebeta.Service
	storageService  *storage.Service
	logger          boshlog.Logger
}

func NewGoogleClient(
	config config.Config,
	logger boshlog.Logger,
) (GoogleClient, error) {
	var err error
	var computeClient, storageClient *http.Client
	userAgent := "bosh-google-cpi/0.0.1"

	if config.JSONKey != "" {
		computeJwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(config.JSONKey), computeScope)
		if err != nil {
			return GoogleClient{}, bosherr.WrapError(err, "Reading Google JSON Key")
		}
		computeClient = computeJwtConf.Client(oauth2.NoContext)

		storageJwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(config.JSONKey), storageScope)
		if err != nil {
			return GoogleClient{}, bosherr.WrapError(err, "Reading Google JSON Key")
		}
		storageClient = storageJwtConf.Client(oauth2.NoContext)
	} else {
		if v := os.Getenv("GCE_METADATA_HOST"); v == "" {
			os.Setenv("GCE_METADATA_HOST", metadataHost)
		}
		computeClient, err = oauthgoogle.DefaultClient(oauth2.NoContext, computeScope)
		if err != nil {
			return GoogleClient{}, bosherr.WrapError(err, "Creating a Google default client")
		}

		storageClient, err = oauthgoogle.DefaultClient(oauth2.NoContext, storageScope)
		if err != nil {
			return GoogleClient{}, bosherr.WrapError(err, "Creating a Google default client")
		}
	}

	// Custom RoundTripper for retries
	computeRetrier := &RetryTransport{
		Base:       computeClient.Transport,
		MaxRetries: 3,
	}
	computeClient.Transport = computeRetrier
	computeService, err := compute.New(computeClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Compute Service client")
	}
	computeService.UserAgent = userAgent

	computeServiceB, err := computebeta.New(computeClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Compute Service client")
	}
	computeServiceB.UserAgent = userAgent

	// Custom RoundTripper for retries
	storageRetrier := &RetryTransport{
		Base:       storageClient.Transport,
		MaxRetries: 3,
	}
	storageClient.Transport = storageRetrier
	storageService, err := storage.New(storageClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Storage Service client")
	}
	storageService.UserAgent = userAgent

	return GoogleClient{
		config:          config,
		computeService:  computeService,
		computeServiceB: computeServiceB,
		storageService:  storageService,
		logger:          logger,
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

func (c GoogleClient) ComputeService() *compute.Service {
	return c.computeService
}

func (c GoogleClient) ComputeBetaService() *computebeta.Service {
	return c.computeServiceB
}

func (c GoogleClient) StorageService() *storage.Service {
	return c.storageService
}
