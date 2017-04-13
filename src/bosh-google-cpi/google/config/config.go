package config

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

var CpiRelease string

type Config struct {
	Project               string `json:"project"`
	UserAgentPrefix       string `json:"user_agent_prefix"`
	JSONKey               string `json:"json_key"`
	DefaultRootDiskSizeGb int    `json:"default_root_disk_size_gb"`
	DefaultRootDiskType   string `json:"default_root_disk_type"`
}

func (c Config) GetUserAgent() string {
	boshCpiUserAgent := "bosh-google-cpi/" + CpiRelease
	if CpiRelease == "" {
		boshCpiUserAgent = boshCpiUserAgent + "dev"
	}
	if c.UserAgentPrefix == "" {
		return boshCpiUserAgent
	}
	return c.UserAgentPrefix + " " + boshCpiUserAgent
}

func (c Config) Validate() error {
	if c.Project == "" {
		return bosherr.Error("Must provide a non-empty Project")
	}
	return nil
}
