package instance

import (
	"fmt"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	computebeta "google.golang.org/api/compute/v0.beta"
)

func (i GoogleInstanceService) SetMetadata(id string, vmMetadata Metadata) error {
	// Find the instance
	instance, found, err := i.FindBeta(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	// We need to reuse the original instance metadata fingerprint and items
	metadata := instance.Metadata
	metadataMap := make(map[string]string)

	// Grab the original metadata items
	for _, item := range metadata.Items {
		metadataMap[item.Key] = *item.Value
	}

	// TODO(evanbrown): Is it possible to update metadata, labels, and tags
	// in a single PATCH request? The current method requires 4 requests to
	// accomplish this.
	// Add or override the new metadata items.
	for key, value := range vmMetadata {
		metadataMap[key] = value
	}

	// Set the new metadata items
	var metadataItems []*computebeta.MetadataItems
	for key, value := range metadataMap {
		mValue := value
		metadataItems = append(metadataItems, &computebeta.MetadataItems{Key: key, Value: &mValue})
	}
	metadata.Items = metadataItems

	i.logger.Debug(googleInstanceServiceLogTag, "Setting metadata for Google Instance '%s'", id)
	operation, err := i.computeServiceB.Instances.SetMetadata(i.project, util.ResourceSplitter(instance.Zone), id, metadata).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to set metadata for Google Instance '%s'", id)
	}

	if _, err = i.operationService.WaiterB(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to set metadata for Google Instance '%s'", id)
	}

	// Apply labels to VM
	// First create a new map and copy existing labels into it
	labelsMap := make(map[string]string)
	for k, v := range instance.Labels {
		labelsMap[k] = v
	}

	for k, v := range vmMetadata {
		safeKey, keyErr := SafeLabel(k, LabelKey)
		safeValue, valueErr := SafeLabel(v, LabelValue)

		if keyErr == nil && valueErr == nil {
			labelsMap[safeKey] = safeValue
		} else {
			i.logger.Debug(googleInstanceServiceLogTag, fmt.Sprintf("Skipped label %q:%q: key error: %v, value error: %v", k, v, keyErr, valueErr))
		}
	}

	labelsRequest := &computebeta.InstancesSetLabelsRequest{
		LabelFingerprint: instance.LabelFingerprint,
		Labels:           labelsMap,
	}
	i.logger.Debug(googleInstanceServiceLogTag, "Setting labels for Google Instance '%s'", id)
	operation, err = i.computeServiceB.Instances.SetLabels(i.project, util.ResourceSplitter(instance.Zone), id, labelsRequest).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to set labels for Google Instance '%s'", id)
	}
	if _, err = i.operationService.WaiterB(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to set labels for Google Instance '%s'", id)
	}

	return nil
}
