package registry

import (
	"encoding/json"
	"fmt"

	"bosh-google-cpi/google/client"
	"bosh-google-cpi/util"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"
)

const (
	GCEMetadataKey       = "bosh_settings"
	metadataClientLogTag = "registryMetadatalient"
)

// MetadataClient represents a GCE metadata client.
type MetadataClient struct {
	options      ClientOptions
	logger       boshlog.Logger
	googleClient client.GoogleClient
}

// NewMetadataClient creates a new GCE metadata client.
func NewMetadataClient(googleClient client.GoogleClient, options ClientOptions, logger boshlog.Logger) MetadataClient {
	if options.GCEMetadataKey == "" {
		options.GCEMetadataKey = GCEMetadataKey
	}
	return MetadataClient{
		googleClient: googleClient,
		options:      options,
		logger:       logger,
	}
}

type instanceMetadata struct {
	items       map[string]string
	fingerprint string
	zone        string
}

func (i *instanceMetadata) computeMetadata() *compute.Metadata {
	metadata := &compute.Metadata{
		Fingerprint: i.fingerprint,
	}

	metadata.Items = make([]*compute.MetadataItems, len(i.items))
	for k := range i.items {
		v := i.items[k]
		metadata.Items = append(metadata.Items, &compute.MetadataItems{Key: k, Value: &v})
	}
	return metadata
}

// Delete deletes the instance settings for a given instance ID.
func (c MetadataClient) Delete(instanceID string) error {
	currentMetadata, err := c.metadata(instanceID)
	if err != nil {
		return err
	}
	delete(currentMetadata.items, c.options.GCEMetadataKey)
	_, err = c.googleClient.ComputeService().Instances.SetMetadata(c.googleClient.Project(), currentMetadata.zone, instanceID, currentMetadata.computeMetadata()).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating instance metadata with SetMetadata call: %v, metadata value: %#v", err, currentMetadata.computeMetadata())
	}

	return nil
}

// Fetch gets the agent settings for a given instance ID.
func (c MetadataClient) Fetch(instanceID string) (AgentSettings, error) {
	metadata, err := c.metadata(instanceID)
	if err != nil {
		return AgentSettings{}, err
	}
	var settings string
	var ok bool
	if settings, ok = metadata.items[c.options.GCEMetadataKey]; !ok {
		return AgentSettings{}, bosherr.WrapError(err, fmt.Sprintf("Retrieving settings from instance metadata for instance %q, key: %q, metadata: %#v", instanceID, c.options.GCEMetadataKey, metadata.items))
	}
	var agentSettings AgentSettings
	if err = json.Unmarshal([]byte(settings), &agentSettings); err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Unmarshalling agent settings response from instance metadata key %q, contents: %v", c.options.GCEMetadataKey, settings)
	}
	return agentSettings, nil

}

// Update updates the agent settings for a given instance ID. If there are not already agent settings for the instance, it will create ones.
func (c MetadataClient) Update(instanceID string, agentSettings AgentSettings) error {
	settingsJSON, err := json.Marshal(agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling agent settings, contents: %#s", agentSettings)
	}
	c.logger.Debug(metadataClientLogTag, "Updating instance metadata for %q with agent settings %q", instanceID, settingsJSON)

	currentMetadata, err := c.metadata(instanceID)
	if err != nil {
		return err
	}
	currentMetadata.items[c.options.GCEMetadataKey] = string(settingsJSON)

	c.logger.Debug(metadataClientLogTag, "Updating instance metadata to: %#v", currentMetadata.computeMetadata())
	_, err = c.googleClient.ComputeService().Instances.SetMetadata(c.googleClient.Project(), currentMetadata.zone, instanceID, currentMetadata.computeMetadata()).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating instance metadata with SetMetadata call: %v, metadata value: %#v", err, currentMetadata.computeMetadata())
	}

	return nil
}

func (c MetadataClient) metadata(instanceID string) (instanceMetadata, error) {
	svc := c.googleClient.ComputeService().Instances
	listCall := svc.AggregatedList(c.googleClient.Project())
	listCall.Filter(fmt.Sprintf("name eq %v", instanceID))
	instances, err := listCall.Do()
	c.logger.Debug(metadataClientLogTag, fmt.Sprintf("Searching for instance %q", instanceID))
	if err != nil {
		return instanceMetadata{}, bosherr.WrapError(err, fmt.Sprintf("Listing instances to find instance %q", instanceID))
	}

	// Find single instance in list
	for _, list := range instances.Items {
		if len(list.Instances) == 1 {
			c.logger.Debug(metadataClientLogTag, fmt.Sprintf("Found instance %q in zone %q", instanceID, list.Instances[0].Zone))
			var metadata instanceMetadata
			metadata.fingerprint = list.Instances[0].Metadata.Fingerprint
			metadata.zone = util.ResourceSplitter(list.Instances[0].Zone)
			metadata.items = make(map[string]string)
			for _, item := range list.Instances[0].Metadata.Items {
				metadata.items[item.Key] = *item.Value
			}
			c.logger.Debug(metadataClientLogTag, fmt.Sprintf("Got metadata for instance %q: %#s", instanceID, metadata))
			return metadata, nil
		}
	}
	return instanceMetadata{}, bosherr.WrapError(err, fmt.Sprintf("Could not find find instance %q", instanceID))
}
