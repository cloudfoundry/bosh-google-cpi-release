package network

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"google.golang.org/api/googleapi"
)

func (n GoogleNetworkService) Find(projectId, id string) (Network, bool, error) {
	n.logger.Debug(googleNetworkServiceLogTag, "Finding Google Network '%s'", id)

	networkItem, err := n.computeService.Networks.Get(n.projectService.Find(projectId), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return Network{}, false, nil
		}

		return Network{}, false, bosherr.WrapErrorf(err, "Failed to find Google Network '%s' in project '%s'", id, projectId)
	}

	network := Network{
		Name:     networkItem.Name,
		SelfLink: networkItem.SelfLink,
	}
	return network, true, nil
}
