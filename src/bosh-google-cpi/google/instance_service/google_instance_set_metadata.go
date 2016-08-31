package instance

import (
	"strings"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	computebeta "google.golang.org/api/compute/v0.beta"
)

// The list of metadata key-value pairs that should be applied as labels
var (
	LabelList []LabelTagMetadata = []LabelTagMetadata{
		{
			Key:       "director",
			CleanerFn: SafeLabel,
		},
		{
			Key:       "name",
			CleanerFn: SafeLabel,
		},
		{
			Key:       "deployment",
			CleanerFn: SafeLabel,
		},
		{
			Key:       "job",
			CleanerFn: SafeLabel,
		},
		{
			Key:       "index",
			CleanerFn: func(s string) string { return "index-" + SafeLabel(s) },
		},
	}

	// The list of metadata keys whose value will be automatically applied as a tag
	TagList []LabelTagMetadata = []LabelTagMetadata{
		{
			Key:       "job",
			CleanerFn: SafeLabel,
		},
		{
			ValFn: func(m map[string]string) string {
				return m["deployment"] + "-" + m["job"]
			},
			CleanerFn: SafeLabel,
		},
		{
			Key:       "deployment",
			CleanerFn: SafeLabel,
		},
	}
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
		metadataMap[item.Key] = item.Value
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

	// Apply labels to VM
	// First create a new map and copy existing labels into it
	labelsMap := make(map[string]string)
	for k, v := range instance.Labels {
		labelsMap[k] = v
	}

	for _, l := range LabelList {
		labelsMap[l.Key] = l.Val(vmMetadata)
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

	// Apply tags to VM
	// Re-retrieve the instance because labels will have changed the tag fingerprint
	instance, _, err = i.FindBeta(id, "")
	if err != nil {
		return err
	}

	// Get existing instance tags
	tags := make(Tags, 0)
	tags = append(tags, Tags(instance.Tags.Items)...)

	// Add metadata specified in TagList to tags
	for _, t := range TagList {
		tags = append(tags, t.Val(vmMetadata))
	}

	// Eliminate duplicate tags
	instance.Tags.Items = tags.Unique()

	i.logger.Debug(googleInstanceServiceLogTag, "Setting tags for Google Instance '%s'", id)
	operation, err = i.computeServiceB.Instances.SetTags(i.project, util.ResourceSplitter(instance.Zone), id, instance.Tags).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to set tags for Google Instance '%s'", id)
	}
	if _, err = i.operationService.WaiterB(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to set tags for Google Instance '%s'", id)
	}

	return nil
}

//
type LabelTagMetadata struct {
	Key       string
	ValFn     func(map[string]string) string
	CleanerFn func(string) string
}

func (l LabelTagMetadata) Val(m map[string]string) string {
	var value string
	if l.ValFn != nil {
		value = l.ValFn(m)
	} else {
		value = m[l.Key]
	}
	return l.CleanerFn(value)
}

func SafeLabel(s string) string {
	maxlen := 63
	// Replace common invalid chars
	s = strings.Replace(s, "/", "-", -1)
	s = strings.Replace(s, "_", "-", -1)

	// Trim to max length
	if len(s) > maxlen {
		s = s[0:maxlen]
	}

	// Ensure the string doesn't end in -
	s = strings.TrimSuffix(s, "-")
	return s
}
