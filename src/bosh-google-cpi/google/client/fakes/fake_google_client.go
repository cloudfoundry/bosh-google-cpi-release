package fakes

import (
	"github.com/cloudfoundry/bosh-utils/logger"

	"bosh-google-cpi/google/client"
	"bosh-google-cpi/google/config"
)

func NewFakeGoogleClient(cfg config.Config, logger logger.Logger) (client.GoogleClient, error) {
	return client.GoogleClient{
		Config: cfg,
	}, nil
}
