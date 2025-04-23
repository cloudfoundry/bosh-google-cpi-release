package instance

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"bosh-google-cpi/util"
)

func (i GoogleInstanceService) Find(id string, zone string) (*compute.Instance, bool, error) {
	if zone == "" {
		i.logger.Debug(googleInstanceServiceLogTag, "Finding Google Instance '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		instances, err := i.computeService.Instances.AggregatedList(i.project).Filter(filter).Do()
		if err != nil {
			return &compute.Instance{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance '%s'", id)
		}

		for _, instanceItems := range instances.Items {
			for _, instance := range instanceItems.Instances {
				// Return the first instance (it can only be 1 instance with the same name across all zones)
				return instance, true, nil
			}
		}

		return &compute.Instance{}, false, nil
	}

	i.logger.Debug(googleInstanceServiceLogTag, "Finding Google Instance '%s' in zone '%s'", id, zone)
	instance, err := i.computeService.Instances.Get(i.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.Instance{}, false, nil
		}

		return &compute.Instance{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance '%s' in zone '%s'", id, util.ResourceSplitter(zone))
	}

	return instance, true, nil
}

func (i GoogleInstanceService) FindBeta(id string, zone string) (*computebeta.Instance, bool, error) {
	if zone == "" {
		i.logger.Debug(googleInstanceServiceLogTag, "Finding Google Instance '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		instances, err := i.computeServiceB.Instances.AggregatedList(i.project).Filter(filter).Do()
		if err != nil {
			return &computebeta.Instance{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance '%s'", id)
		}

		for _, instanceItems := range instances.Items {
			for _, instance := range instanceItems.Instances {
				// Return the first instance (it can only be 1 instance with the same name across all zones)
				return instance, true, nil
			}
		}

		return &computebeta.Instance{}, false, nil
	}

	i.logger.Debug(googleInstanceServiceLogTag, "Finding Google Instance '%s' in zone '%s'", id, zone)
	instance, err := i.computeServiceB.Instances.Get(i.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &computebeta.Instance{}, false, nil
		}

		return &computebeta.Instance{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance '%s' in zone '%s'", id, util.ResourceSplitter(zone))
	}

	return instance, true, nil
}
