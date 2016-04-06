package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bgcaction "bosh-google-cpi/action"
	bogcconfig "bosh-google-cpi/google/config"
)

type Config struct {
	Google bogcconfig.Config

	Actions bgcaction.ConcreteFactoryOptions
}

func NewConfigFromPath(configFile string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	if configFile == "" {
		return config, bosherr.Errorf("Must provide a config file")
	}

	bytes, err := fs.ReadFile(configFile)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config file '%s'", configFile)
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config contents")
	}

	if err = config.Validate(); err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

func (c Config) Validate() error {
	if err := c.Google.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating Google configuration")
	}

	if err := c.Actions.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating Actions configuration")
	}

	return nil
}
