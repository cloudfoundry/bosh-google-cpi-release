package store

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	"github.com/mitchellh/mapstructure"
)

type RegistryStore interface {
	Delete(string) error
	Get(string) (string, bool, error)
	Save(string, string) error
}

func NewRegistryStore(
	config RegistryStoreConfig,
	logger boshlog.Logger,
) (RegistryStore, error) {
	switch {
	case config.Adapter == "bolt":
		boltRegistryStoreConfig := BoltRegistryStoreConfig{}
		err := mapstructure.Decode(config.Options, &boltRegistryStoreConfig)
		if err != nil {
			return nil, bosherr.WrapError(err, "Decoding Bolt Registry Store configuration")
		}

		err = boltRegistryStoreConfig.Validate()
		if err != nil {
			return nil, bosherr.WrapError(err, "Validating Bolt Registry Store configuration")
		}

		return NewBoltRegistryStore(boltRegistryStoreConfig, logger), nil
	}

	return nil, bosherr.Errorf("Registry Store adapter '%s' not supported", config.Adapter)
}
