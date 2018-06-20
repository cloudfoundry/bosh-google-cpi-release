package fakes

import (
	"bosh-google-cpi/google/client"
	"bosh-google-cpi/google/config"
	"github.com/cloudfoundry/bosh-utils/logger"
)

func NewFakeGoogleClient(cfg config.Config, logger logger.Logger) (client.GoogleClient, error) {
	return client.GoogleClient{
		Config: cfg,
	}, nil
}
