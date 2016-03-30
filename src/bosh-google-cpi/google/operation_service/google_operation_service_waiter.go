package operation

import (
	"math"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (o GoogleOperationService) Waiter(operation *compute.Operation, zone string, region string) (*compute.Operation, error) {
	var tries int
	var err error

	start := time.Now()
	for tries = 1; tries < googleOperationServiceMaxTries; tries++ {
		factor := math.Pow(2, math.Min(float64(tries), float64(googleOperationServiceMaxSleepExponent)))
		wait := time.Duration(factor) * time.Second
		o.logger.Debug(googleOperationServiceLogTag, "Waiting for Google Operation '%s' to be ready, retrying in %v (%d/%d)", operation.Name, wait, tries, googleOperationServiceMaxTries)
		time.Sleep(wait)

		if zone == "" {
			if region == "" {
				operation, err = o.computeService.GlobalOperations.Get(o.project, operation.Name).Do()
			} else {
				operation, err = o.computeService.RegionOperations.Get(o.project, util.ResourceSplitter(region), operation.Name).Do()
			}
		} else {
			operation, err = o.computeService.ZoneOperations.Get(o.project, util.ResourceSplitter(zone), operation.Name).Do()
		}

		if err != nil {
			o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' finished with an error: %#v", operation.Name, err)
			if operation.Error != nil {
				return nil, bosherr.WrapErrorf(GoogleOperationError(*operation.Error), "Google Operation '%s' finished with an error", operation.Name)
			}

			return nil, bosherr.WrapErrorf(err, "Google Operation '%s' finished with an error", operation.Name)
		}

		if operation.Status == googleOperationReadyStatus {
			if operation.Error != nil {
				o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' finished with an error: %s", operation.Name, GoogleOperationError(*operation.Error))
				return nil, bosherr.WrapErrorf(GoogleOperationError(*operation.Error), "Google Operation '%s' finished with an error", operation.Name)
			}

			o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' is now ready after %v", operation.Name, time.Since(start))
			return operation, nil
		}
	}

	return nil, bosherr.Errorf("Timed out waiting for Google Operation '%s' to be ready", operation.Name)
}
