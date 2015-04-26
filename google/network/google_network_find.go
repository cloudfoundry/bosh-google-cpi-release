package gnetwork

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (n GoogleNetworkService) Find(id string) (*compute.Network, bool, error) {
	n.logger.Debug(googleNetworkServiceLogTag, "Finding Google Network '%s'", id)
	network, err := n.computeService.Networks.Get(n.project, id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.Network{}, false, nil
		}

		return &compute.Network{}, false, bosherr.WrapErrorf(err, "Failed to find Google Network '%s'", id)
	}

	return network, true, nil
}
