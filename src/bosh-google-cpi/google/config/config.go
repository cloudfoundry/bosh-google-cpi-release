package config

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Config struct {
	Project               string `json:"project"`
	UserAgent             string `json:"user_agent"`
	JSONKey               string `json:"json_key"`
	DefaultRootDiskSizeGb int    `json:"default_root_disk_size_gb"`
	DefaultRootDiskType   string `json:"default_root_disk_type"`
}

func (c Config) GetUserAgent() string {
	boshCpiUserAgent := "bosh-google-cpi/0.0.1"
	if c.UserAgent == "" {
		return boshCpiUserAgent
	}
	return c.UserAgent + " " + boshCpiUserAgent
}

func (c Config) Validate() error {
	if c.Project == "" {
		return bosherr.Error("Must provide a non-empty Project")
	}
	return nil
}
