package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

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

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config %s", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	err = config.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

func (c Config) Validate() error {
	err := c.Google.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Google configuration")
	}

	err = c.Actions.Validate()
	if err != nil {
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
