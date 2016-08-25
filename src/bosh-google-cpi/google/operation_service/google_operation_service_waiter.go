package operation

import (
	"math"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func (o GoogleOperationService) Waiter(operation *compute.Operation, zone string, region string) (*compute.Operation, error) {
	var tries int
	var err error
	var opName string

	start := time.Now()
	for tries = 1; tries < googleOperationServiceMaxTries; tries++ {
		factor := math.Pow(2, math.Min(float64(tries), float64(googleOperationServiceMaxSleepExponent)))
		wait := time.Duration(factor) * time.Second
		opName = operation.Name
		o.logger.Debug(googleOperationServiceLogTag, "Waiting for Google Operation '%s' to be ready, retrying in %v (%d/%d)", opName, wait, tries, googleOperationServiceMaxTries)
		time.Sleep(wait)

		if zone == "" {
			if region == "" {
				operation, err = o.computeService.GlobalOperations.Get(o.project, opName).Do()
			} else {
				operation, err = o.computeService.RegionOperations.Get(o.project, util.ResourceSplitter(region), opName).Do()
			}
		} else {
			operation, err = o.computeService.ZoneOperations.Get(o.project, util.ResourceSplitter(zone), opName).Do()
		}

		if err != nil {
			o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' finished with an error: %#v", opName, err)
			if operation != nil && operation.Error != nil {
				return nil, bosherr.WrapErrorf(GoogleOperationError(*operation.Error), "Google Operation '%s' finished with an error", opName)
			}

			return nil, bosherr.WrapErrorf(err, "Google Operation '%s' finished with an error", opName)
		}

		if operation.Status == googleOperationReadyStatus {
			if operation.Error != nil {
				o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' finished with an error: %s", opName, GoogleOperationError(*operation.Error))
				return nil, bosherr.WrapErrorf(GoogleOperationError(*operation.Error), "Google Operation '%s' finished with an error", opName)
			}

			o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' is now ready after %v", opName, time.Since(start))
			return operation, nil
		}
	}

	return nil, bosherr.Errorf("Timed out waiting for Google Operation '%s' to be ready", opName)
}

func (o GoogleOperationService) WaiterB(operation *computebeta.Operation, zone string, region string) (*computebeta.Operation, error) {
	var tries int
	var err error
	var opName string

	start := time.Now()
	for tries = 1; tries < googleOperationServiceMaxTries; tries++ {
		factor := math.Pow(2, math.Min(float64(tries), float64(googleOperationServiceMaxSleepExponent)))
		wait := time.Duration(factor) * time.Second
		opName = operation.Name
		o.logger.Debug(googleOperationServiceLogTag, "Waiting for Google Operation '%s' to be ready, retrying in %v (%d/%d)", opName, wait, tries, googleOperationServiceMaxTries)
		time.Sleep(wait)

		if zone == "" {
			if region == "" {
				operation, err = o.computeServiceB.GlobalOperations.Get(o.project, opName).Do()
			} else {
				operation, err = o.computeServiceB.RegionOperations.Get(o.project, util.ResourceSplitter(region), opName).Do()
			}
		} else {
			operation, err = o.computeServiceB.ZoneOperations.Get(o.project, util.ResourceSplitter(zone), opName).Do()
		}

		if err != nil {
			o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' finished with an error: %#v", opName, err)
			if operation != nil && operation.Error != nil {
				return nil, bosherr.WrapErrorf(GoogleOperationErrorB(*operation.Error), "Google Operation '%s' finished with an error", opName)
			}

			return nil, bosherr.WrapErrorf(err, "Google Operation '%s' finished with an error", opName)
		}

		if operation.Status == googleOperationReadyStatus {
			if operation.Error != nil {
				o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' finished with an error: %s", opName, GoogleOperationErrorB(*operation.Error))
				return nil, bosherr.WrapErrorf(GoogleOperationErrorB(*operation.Error), "Google Operation '%s' finished with an error", opName)
			}

			o.logger.Debug(googleOperationServiceLogTag, "Google Operation '%s' is now ready after %v", opName, time.Since(start))
			return operation, nil
		}
	}

	return nil, bosherr.Errorf("Timed out waiting for Google Operation '%s' to be ready", opName)
}
