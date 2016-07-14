package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/registry"
)

type ConcreteFactoryOptions struct {
	Agent    registry.AgentOptions
	Registry registry.ClientOptions

	// List of NTP servers
	Ntp []string

	// Blobstore options
	Blobstore registry.BlobstoreOptions
}

func (o ConcreteFactoryOptions) Validate() error {
	if err := o.Agent.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating Agent configuration")
	}
	err := o.Blobstore.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Blobstore configuration")
	}
	if err := o.Registry.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating Registry configuration")
	}

	return nil
}
