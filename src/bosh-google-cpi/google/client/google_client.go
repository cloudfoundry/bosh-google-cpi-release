package client

import (
	"net/http"
	"os"
	"time"

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

	// Configuration for retrier.
	retries         = 12
	firstRetrySleep = 50 * time.Millisecond
)

type GoogleClient struct {
	Config          config.Config
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
	var cloudClient *http.Client
	userAgent := config.GetUserAgent()

	if config.JSONKey != "" {
		jwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(config.JSONKey), computeScope, storageScope)
		if err != nil {
			return GoogleClient{}, bosherr.WrapError(err, "Reading Google JSON Key")
		}
		cloudClient = jwtConf.Client(oauth2.NoContext)
	} else {
		if v := os.Getenv("GCE_METADATA_HOST"); v == "" {
			os.Setenv("GCE_METADATA_HOST", metadataHost)
		}
		cloudClient, err = oauthgoogle.DefaultClient(oauth2.NoContext, computeScope, storageScope)
		if err != nil {
			return GoogleClient{}, bosherr.WrapError(err, "Creating a Google default client")
		}
	}

	// Custom RoundTripper for retries
	retrier := &RetryTransport{
		Base:            cloudClient.Transport,
		MaxRetries:      retries,
		FirstRetrySleep: firstRetrySleep,
		logger:          logger,
	}
	cloudClient.Transport = retrier

	computeService, err := compute.New(cloudClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Compute Service client")
	}
	computeService.UserAgent = userAgent

	computeServiceB, err := computebeta.New(cloudClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Compute Service client")
	}
	computeServiceB.UserAgent = userAgent

	storageService, err := storage.New(cloudClient)
	if err != nil {
		return GoogleClient{}, bosherr.WrapError(err, "Creating a Google Storage Service client")
	}
	storageService.UserAgent = userAgent

	return GoogleClient{
		Config:          config,
		computeService:  computeService,
		computeServiceB: computeServiceB,
		storageService:  storageService,
		logger:          logger,
	}, nil
}

func (c GoogleClient) Project() string {
	return c.Config.Project
}

func (c GoogleClient) DefaultRootDiskSizeGb() int {
	return c.Config.DefaultRootDiskSizeGb
}

func (c GoogleClient) DefaultRootDiskType() string {
	return c.Config.DefaultRootDiskType
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
