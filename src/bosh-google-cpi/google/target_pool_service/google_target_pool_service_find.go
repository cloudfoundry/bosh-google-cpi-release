package targetpool

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/googleapi"
)

func (t GoogleTargetPoolService) Find(id string, region string) (TargetPool, bool, error) {
	if region == "" {
		t.logger.Debug(googleTargetPoolServiceLogTag, "Finding Google Target Pool '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		targetPools, err := t.computeService.TargetPools.AggregatedList(t.project).Filter(filter).Do()
		if err != nil {
			return TargetPool{}, false, bosherr.WrapErrorf(err, "Failed to find Google Target Pool '%s'", id)
		}

		for _, targetPoolItems := range targetPools.Items {
			for _, targetPoolItem := range targetPoolItems.TargetPools {
				// Return the first target pool (it can only be 1 target pool with the same name across all regions)
				targetPool := TargetPool{
					Name:      targetPoolItem.Name,
					Instances: targetPoolItem.Instances,
					SelfLink:  targetPoolItem.SelfLink,
					Region:    targetPoolItem.Region,
				}
				return targetPool, true, nil
			}
		}

		return TargetPool{}, false, nil
	}

	t.logger.Debug(googleTargetPoolServiceLogTag, "Finding Google Target Pool '%s' in region '%s'", id, region)
	targetPoolItem, err := t.computeService.TargetPools.Get(t.project, util.ResourceSplitter(region), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return TargetPool{}, false, nil
		}

		return TargetPool{}, false, bosherr.WrapErrorf(err, "Failed to find Google Target Pool '%s' in region '%s'", id, region)
	}

	targetPool := TargetPool{
		Name:      targetPoolItem.Name,
		Instances: targetPoolItem.Instances,
		SelfLink:  targetPoolItem.SelfLink,
		Region:    targetPoolItem.Region,
	}
	return targetPool, true, nil
}
