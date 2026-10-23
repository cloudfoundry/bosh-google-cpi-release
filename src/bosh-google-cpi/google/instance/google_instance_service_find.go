package instance

import (
	"context"
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
		filter := fmt.Sprintf("name = %q", id)
		var found *compute.Instance
		err := i.computeService.Instances.AggregatedList(i.project).Filter(filter).Pages(context.Background(), func(page *compute.InstanceAggregatedList) error {
			for _, items := range page.Items {
				for _, instance := range items.Instances {
					if instance.Name != id {
						continue
					}
					if found != nil {
						return fmt.Errorf("instance %q is ambiguous across zones", id)
					}
					found = instance
				}
			}
			return nil
		})
		if err != nil {
			return &compute.Instance{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance '%s'", id)
		}
		if found != nil {
			return found, true, nil
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
		filter := fmt.Sprintf("name = %q", id)
		var found *computebeta.Instance
		err := i.computeServiceB.Instances.AggregatedList(i.project).Filter(filter).Pages(context.Background(), func(page *computebeta.InstanceAggregatedList) error {
			for _, items := range page.Items {
				for _, instance := range items.Instances {
					if instance.Name != id {
						continue
					}
					if found != nil {
						return fmt.Errorf("instance %q is ambiguous across zones", id)
					}
					found = instance
				}
			}
			return nil
		})
		if err != nil {
			return &computebeta.Instance{}, false, bosherr.WrapErrorf(err, "Failed to find Google Instance '%s'", id)
		}
		if found != nil {
			return found, true, nil
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
