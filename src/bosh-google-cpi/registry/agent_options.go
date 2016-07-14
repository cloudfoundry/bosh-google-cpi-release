package registry

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

// AgentOptions are the agent options passed to the the BOSH Agent (http://bosh.io/docs/bosh-components.html#agent).
type AgentOptions struct {
	// Mbus URI
	Mbus string
}

// BlobstoreOptions are the blobstore options passed to the BOSH Agent (http://bosh.io/docs/bosh-components.html#agent).
type BlobstoreOptions struct {
	// Blobstore type
	Type string

	// Blobstore options
	Options map[string]interface{}
}

// Validate validates the Agent options.
func (o AgentOptions) Validate() error {
	if o.Mbus == "" {
		return bosherr.Error("Must provide a non-empty Mbus")
	}

	return nil
}

// Validate validates the Blobstore options.
func (o BlobstoreOptions) Validate() error {
	if o.Type == "" {
		return bosherr.Error("Must provide non-empty Type")
	}

	return nil
}
