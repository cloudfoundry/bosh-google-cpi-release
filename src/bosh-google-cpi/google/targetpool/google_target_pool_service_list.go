package targetpool

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
)

func (t GoogleTargetPoolService) List(region string) ([]TargetPool, error) {
	var targetPools []TargetPool

	if region == "" {
		t.logger.Debug(googleTargetPoolServiceLogTag, "Listing Google Target Pools")
		targetPoolAggregatedList, err := t.computeService.TargetPools.AggregatedList(t.project).Do()
		if err != nil {
			return targetPools, bosherr.WrapError(err, "Failed to list Google Target Pools")
		}

		for _, targetPoolList := range targetPoolAggregatedList.Items {
			for _, targetPoolItem := range targetPoolList.TargetPools {
				targetPool := TargetPool{
					Name:      targetPoolItem.Name,
					Instances: targetPoolItem.Instances,
					SelfLink:  targetPoolItem.SelfLink,
					Region:    targetPoolItem.Region,
				}
				targetPools = append(targetPools, targetPool)
			}
		}

		return targetPools, nil
	}

	t.logger.Debug(googleTargetPoolServiceLogTag, "Listing Google Target Pools in region '%s'", region)
	targetPoolList, err := t.computeService.TargetPools.List(t.project, util.ResourceSplitter(region)).Do()
	if err != nil {
		return targetPools, bosherr.WrapErrorf(err, "Failed to list Google Target Pools in region '%s'", region)
	}

	for _, targetPoolItem := range targetPoolList.Items {
		targetPool := TargetPool{
			Name:      targetPoolItem.Name,
			Instances: targetPoolItem.Instances,
			SelfLink:  targetPoolItem.SelfLink,
			Region:    targetPoolItem.Region,
		}
		targetPools = append(targetPools, targetPool)
	}

	return targetPools, nil
}
