package subnetwork

import (
	"errors"

	"bosh-google-cpi/util"
)

var ErrRegionRequired error = errors.New("A region is required to find a subnet")
var ErrSubnetNotFound error = errors.New("Subnet could not be found")

func (s GoogleSubnetworkService) Find(projectId, id, region string) (Subnetwork, error) {
	if region == "" {
		return Subnetwork{}, ErrRegionRequired
	}

	// Default to compute project if not specified
	if projectId == "" {
		projectId = s.project
	}

	s.logger.Debug(googleSubnetworkServiceLogTag, "Finding Google Subnetwork '%s' in region '%s' in project '%s'", id, region, projectId)
	subnetworkItem, err := s.computeService.Subnetworks.Get(projectId, util.ResourceSplitter(region), id).Do()
	if err != nil {
		return Subnetwork{}, err
	}

	subnetwork := Subnetwork{
		Name:     subnetworkItem.Name,
		SelfLink: subnetworkItem.SelfLink,
	}
	return subnetwork, nil
}
