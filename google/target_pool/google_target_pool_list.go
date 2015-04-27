package gtargetpool

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
)

func (t GoogleTargetPoolService) List(region string) ([]*compute.TargetPool, error) {
	var targetPools []*compute.TargetPool

	if region == "" {
		t.logger.Debug(googleTargetPoolServiceLogTag, "Listing Google Target Pools")
		targetPoolAggregatedList, err := t.computeService.TargetPools.AggregatedList(t.project).Do()
		if err != nil {
			return targetPools, bosherr.WrapError(err, "Failed to list Google Target Pools")
		}

		for _, targetPoolList := range targetPoolAggregatedList.Items {
			for _, targetPool := range targetPoolList.TargetPools {
				targetPools = append(targetPools, targetPool)
			}
		}

		return targetPools, nil
	}

	t.logger.Debug(googleTargetPoolServiceLogTag, "Listing Google Target Pools in region '%s'", region)
	targetPoolList, err := t.computeService.TargetPools.List(t.project, gutil.ResourceSplitter(region)).Do()
	if err != nil {
		return targetPools, bosherr.WrapErrorf(err, "Failed to list Google Target Pools in region '%s'", region)
	}

	for _, targetPool := range targetPoolList.Items {
		targetPools = append(targetPools, targetPool)
	}

	return targetPools, nil
}
