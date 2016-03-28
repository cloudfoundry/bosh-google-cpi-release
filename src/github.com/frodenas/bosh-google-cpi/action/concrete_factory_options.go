package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-registry/client"
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
