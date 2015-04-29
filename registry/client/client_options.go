package registry

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

type ClientOptions struct {
	Schema   string
	Host     string
	Port     int
	Username string
	Password string
}

func (o ClientOptions) Validate() error {
	if o.Schema == "" {
		return bosherr.Error("Must provide a non-empty Schema")
	}

	if o.Host == "" {
		return bosherr.Error("Must provide a non-empty Host")
	}

	if o.Port == 0 {
		return bosherr.Error("Must provide a non-empty Port")
	}

	if o.Username == "" {
		return bosherr.Error("Must provide a non-empty Username")
	}

	if o.Password == "" {
		return bosherr.Error("Must provide a non-empty Password")
	}

	return nil
}
