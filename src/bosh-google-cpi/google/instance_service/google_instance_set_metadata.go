package instance

import (
	"strings"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	computebeta "google.golang.org/api/compute/v0.beta"
)

// The list of metadata key-value pairs that should be applied as labels
var labelList []string = []string{"director", "name", "id", "deployment", "job"}

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
		metadataMap[item.Key] = item.Value
	}

	// Add or override the new metadata items.
	for key, value := range vmMetadata {
		metadataMap[key] = value.(string)
	}

	// Set the new metadata items
	var metadataItems []*computebeta.MetadataItems
	for key, value := range metadataMap {
		mValue := value
		metadataItems = append(metadataItems, &computebeta.MetadataItems{Key: key, Value: mValue})
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

	// Repeat the metadata process, but with labels
	labelsMap := make(map[string]string)
	for k, v := range instance.Labels {
		labelsMap[k] = v
	}
	for _, l := range labelList {
		if v, ok := vmMetadata[l]; ok {
			labelsMap[l] = saveLabel(v.(string))
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

func safeLabel(s string, maxlen int) string {
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, "-", "", -1)
	if len(s) > maxlen {
		s = s[0:maxlen]
	}
	return s
}
