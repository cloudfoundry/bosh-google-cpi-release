package subnetwork

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/googleapi"
)

func (s GoogleSubnetworkService) Find(id string, region string) (Subnetwork, bool, error) {
	if region == "" {
		s.logger.Debug(googleSubnetworkServiceLogTag, "Finding Google Subnetwork '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		subnetworks, err := s.computeService.Subnetworks.AggregatedList(s.project).Filter(filter).Do()
		if err != nil {
			return Subnetwork{}, false, bosherr.WrapErrorf(err, "Failed to find Google Subnetwork '%s'", id)
		}

		for _, subnetworkItems := range subnetworks.Items {
			for _, subnetworkItem := range subnetworkItems.Subnetworks {
				// Return the first subnetwork (it can only be 1 subnetwork with the same name across all regions)
				subnetwork := Subnetwork{
					Name:     subnetworkItem.Name,
					SelfLink: subnetworkItem.SelfLink,
				}
				return subnetwork, true, nil
			}
		}

		return Subnetwork{}, false, nil
	}

	s.logger.Debug(googleSubnetworkServiceLogTag, "Finding Google Subnetwork '%s' in region '%s'", id, region)
	subnetworkItem, err := s.computeService.Subnetworks.Get(s.project, util.ResourceSplitter(region), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return Subnetwork{}, false, nil
		}

		return Subnetwork{}, false, bosherr.WrapErrorf(err, "Failed to find Google Subnetwork '%s' in region '%s'", id, region)
	}

	subnetwork := Subnetwork{
		Name:     subnetworkItem.Name,
		SelfLink: subnetworkItem.SelfLink,
	}
	return subnetwork, true, nil
}
