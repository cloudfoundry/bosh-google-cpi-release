package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bgcaction "github.com/frodenas/bosh-google-cpi/action"
)

type Config struct {
	Google GoogleConfig

	Actions bgcaction.ConcreteFactoryOptions
}

type GoogleConfig struct {
	Project     string `json:"project"`
	JSONKey     string `json:"json_key"`
	DefaultZone string `json:"default_zone"`
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

func (c GoogleConfig) Validate() error {
	if c.Project == "" {
		return bosherr.Error("Must provide a non-empty Project")
	}

	if c.JSONKey == "" {
		return bosherr.Error("Must provide a non-empty JSONKey")
	}

	if c.DefaultZone == "" {
		return bosherr.Error("Must provide a non-empty DefaultZone")
	}

	return nil
}
