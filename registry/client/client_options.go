package registry

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

type ClientOptions struct {
	Schema   string    `json:"schema,omitempty"`
	Host     string    `json:"host,omitempty"`
	Port     int       `json:"port,omitempty"`
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	TLS      TLSConfig `json:"tls,omitempty"`
}

type TLSConfig struct {
	InsecureSkipVerify bool   `json:"insecure_skip_verify,omitempty"`
	CertFile           string `json:"certfile,omitempty"`
	KeyFile            string `json:"keyfile,omitempty"`
	CACertFile         string `json:"cacertfile,omitempty"`
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

	if o.Schema == "https" {
		err := o.TLS.Validate()
		if err != nil {
			return bosherr.WrapError(err, "Validating TLS configuration")
		}
	}

	return nil
}

func (o TLSConfig) Validate() error {
	if o.CertFile == "" {
		return bosherr.Error("Must provide a non-empty CertFile")
	}

	if o.KeyFile == "" {
		return bosherr.Error("Must provide a non-empty KeyFile")
	}

	return nil
}
