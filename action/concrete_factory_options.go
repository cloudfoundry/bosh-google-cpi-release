package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/registry/client"
)

type ConcreteFactoryOptions struct {
	Agent    registry.AgentOptions
	Registry registry.ClientOptions
}

func (o ConcreteFactoryOptions) Validate() error {
	err := o.Agent.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Agent configuration")
	}

	err = o.Registry.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Registry configuration")
	}

	return nil
}
