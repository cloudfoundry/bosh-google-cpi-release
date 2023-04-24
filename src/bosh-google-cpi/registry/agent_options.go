package registry

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

// AgentOptions are the agent options passed to the the BOSH Agent (http://bosh.io/docs/bosh-components.html#agent).
type AgentOptions struct {
	// Mbus URI
	Mbus string

	// List of NTP servers
	Ntp []string
}

// Validate validates the Agent options.
func (o AgentOptions) Validate() error {
	if o.Mbus == "" {
		return bosherr.Error("Must provide a non-empty Mbus")
	}

	return nil
}
