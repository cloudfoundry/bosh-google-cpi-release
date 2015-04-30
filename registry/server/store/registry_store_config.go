package store

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

type RegistryStoreConfig struct {
	Adapter string                 `json:"adapter,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

func (c RegistryStoreConfig) Validate() error {
	if c.Adapter == "" {
		return bosherr.Error("Must provide a non-empty Adapter")
	}

	return nil
}
