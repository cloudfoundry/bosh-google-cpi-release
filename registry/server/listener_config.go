package server

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

type ListenerConfig struct {
	Protocol string    `json:"protocol,omitempty"`
	Address  string    `json:"address,omitempty"`
	Port     int       `json:"port,omitempty"`
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	TLS      TLSConfig `json:"tls,omitempty"`
}

type TLSConfig struct {
	CertFile   string `json:"certfile,omitempty"`
	KeyFile    string `json:"keyfile,omitempty"`
	CACertFile string `json:"cacertfile,omitempty"`
}

func (c ListenerConfig) Validate() error {
	if c.Protocol != "http" && c.Protocol != "https" {
		return bosherr.Error("Must provide a valid Protocol")
	}

	if c.Address == "" {
		return bosherr.Error("Must provide a non-empty Address")
	}

	if c.Port == 0 {
		return bosherr.Error("Must provide a non-empty Port")
	}

	if c.Username == "" {
		return bosherr.Error("Must provide a non-empty Username")
	}

	if c.Password == "" {
		return bosherr.Error("Must provide a non-empty Password")
	}

	if c.Protocol == "https" {
		err := c.TLS.Validate()
		if err != nil {
			return bosherr.WrapError(err, "Validating TLS configuration")
		}
	}

	return nil
}

func (c TLSConfig) Validate() error {
	if c.CertFile == "" {
		return bosherr.Error("Must provide a non-empty CertFile")
	}

	if c.KeyFile == "" {
		return bosherr.Error("Must provide a non-empty KeyFile")
	}

	return nil
}
