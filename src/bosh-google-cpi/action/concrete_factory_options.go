package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/registry"
)

type ConcreteFactoryOptions struct {
	Agent    registry.AgentOptions
	Registry registry.ClientOptions
}

func (o ConcreteFactoryOptions) Validate() error {
	if err := o.Agent.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating Agent configuration")
	}

	if err := o.Registry.Validate(); err != nil {
		return bosherr.WrapError(err, "Validating Registry configuration")
	}

	return nil
}
