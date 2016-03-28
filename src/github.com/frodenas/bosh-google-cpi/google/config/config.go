package config

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Config struct {
	Project               string `json:"project"`
	DefaultZone           string `json:"default_zone"`
	JSONKey               string `json:"json_key"`
	DefaultRootDiskSizeGb int    `json:"default_root_disk_size_gb"`
	DefaultRootDiskType   string `json:"default_root_disk_type"`
}

func (c Config) Validate() error {
	if c.Project == "" {
		return bosherr.Error("Must provide a non-empty Project")
	}

	if c.DefaultZone == "" {
		return bosherr.Error("Must provide a non-empty DefaultZone")
	}

	return nil
}
