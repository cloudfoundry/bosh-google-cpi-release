package store

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

type BoltRegistryStoreConfig struct {
	DBFile string
}

func (c BoltRegistryStoreConfig) Validate() error {
	if c.DBFile == "" {
		return bosherr.Error("Must provide a non-empty DBFile")
	}

	return nil
}
