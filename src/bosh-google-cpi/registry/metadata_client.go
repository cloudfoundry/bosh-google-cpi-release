package registry

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"bosh-google-cpi/google/client"
	opsvc "bosh-google-cpi/google/operation_service"
	"bosh-google-cpi/util"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"
)

const (
	GCEMetadataKey       = "bosh_settings"
	metadataClientLogTag = "registryMetadataClient"
	opWaiterRetryMax     = 100
	opMaxSleepExponent   = 3
	opReadyStatus        = "DONE"
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

// Delete deletes the instance settings for a given instance ID. This
// is a no-op as metadata is associated with the instance and is
// deleted when the VM is terminated. There is no need to clean up
// metadata after a VM is deleted.
func (c MetadataClient) Delete(instanceID string) error {
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
		return bosherr.WrapErrorf(err, "Marshalling agent settings, contents: %#v", agentSettings)
	}
	c.logger.Debug(metadataClientLogTag, "Updating instance metadata for %q with agent settings %q", instanceID, settingsJSON)

	currentMetadata, err := c.metadata(instanceID)
	if err != nil {
		return err
	}
	currentMetadata.items[c.options.GCEMetadataKey] = string(settingsJSON)

	computedMetadata := currentMetadata.computeMetadata()

	c.logger.Debug(metadataClientLogTag, "Updating instance metadata to: %#v", computedMetadata)
	op, err := c.googleClient.ComputeService().Instances.SetMetadata(c.googleClient.Project(), currentMetadata.zone, instanceID, computedMetadata).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating instance metadata with SetMetadata call: %v, metadata value: %#v", err, computedMetadata)
	}
	_, err = c.wait(op)
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating instance metadata with SetMetadata call: %v, metadata value: %#v", err, computedMetadata)
	}
	return nil
}

func (c MetadataClient) metadata(instanceID string) (instanceMetadata, error) {
	svc := c.googleClient.ComputeService().Instances
	listCall := svc.AggregatedList(c.googleClient.Project())
	listCall.Filter(fmt.Sprintf("name eq %v", instanceID))
	aggregatedList, err := listCall.Do()
	c.logger.Debug(metadataClientLogTag, fmt.Sprintf("Searching for instance %q", instanceID))
	if err != nil {
		return instanceMetadata{}, bosherr.WrapError(err, fmt.Sprintf("Listing instances to find instance %q", instanceID))
	}

	// Find single instance in list
	for _, list := range aggregatedList.Items {
		if len(list.Instances) == 1 {
			c.logger.Debug(metadataClientLogTag, fmt.Sprintf("Found instance %q in zone %q", instanceID, list.Instances[0].Zone))
			var metadata instanceMetadata
			metadata.fingerprint = list.Instances[0].Metadata.Fingerprint
			metadata.zone = util.ResourceSplitter(list.Instances[0].Zone)
			metadata.items = make(map[string]string)
			for _, item := range list.Instances[0].Metadata.Items {
				metadata.items[item.Key] = *item.Value
			}
			c.logger.Debug(metadataClientLogTag, fmt.Sprintf("Got metadata for instance %q: %#v", instanceID, metadata))
			return metadata, nil
		}
	}
	return instanceMetadata{}, bosherr.WrapError(err, fmt.Sprintf("Could not find find instance %q", instanceID))
}

func (c MetadataClient) wait(operation *compute.Operation) (*compute.Operation, error) {
	var tries int
	var err error

	start := time.Now()
	for tries = 1; tries < opWaiterRetryMax; tries++ {
		factor := math.Pow(2, math.Min(float64(tries), float64(opMaxSleepExponent)))
		wait := time.Duration(factor) * time.Second
		c.logger.Debug(metadataClientLogTag, "Waiting for Google Operation '%s' to be ready, retrying in %v (%d/%d)", operation.Name, wait, tries, opWaiterRetryMax)
		time.Sleep(wait)

		operation, err = c.googleClient.ComputeService().ZoneOperations.Get(c.googleClient.Project(), util.ResourceSplitter(operation.Zone), operation.Name).Do()

		if err != nil {
			c.logger.Debug(metadataClientLogTag, "Google Operation '%s' finished with an error: %#v", operation.Name, err)
			if operation.Error != nil {
				return nil, bosherr.WrapErrorf(opsvc.GoogleOperationError(*operation.Error), "Google Operation '%s' finished with an error", operation.Name)
			}

			return nil, bosherr.WrapErrorf(err, "Google Operation '%s' finished with an error", operation.Name)
		}

		if operation.Status == opReadyStatus {
			if operation.Error != nil {
				c.logger.Debug(metadataClientLogTag, "Google Operation '%s' finished with an error: %s", operation.Name, opsvc.GoogleOperationError(*operation.Error))
				return nil, bosherr.WrapErrorf(opsvc.GoogleOperationError(*operation.Error), "Google Operation '%s' finished with an error", operation.Name)
			}

			c.logger.Debug(metadataClientLogTag, "Google Operation '%s' is now ready after %v", operation.Name, time.Since(start))
			return operation, nil
		}
	}

	return nil, bosherr.Errorf("Timed out waiting for Google Operation '%s' to be ready", operation.Name)
}
