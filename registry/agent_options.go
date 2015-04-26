package registry

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

type AgentOptions struct {
	Mbus      string
	Ntp       []string
	Blobstore BlobstoreOptions
}

type BlobstoreOptions struct {
	Type    string
	Options map[string]interface{}
}

func (o AgentOptions) Validate() error {
	if o.Mbus == "" {
		return bosherr.Error("Must provide a non-empty Mbus")
	}

	err := o.Blobstore.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Blobstore configuration")
	}

	return nil
}

func (o BlobstoreOptions) Validate() error {
	if o.Type == "" {
		return bosherr.Error("Must provide non-empty Type")
	}

	return nil
}
