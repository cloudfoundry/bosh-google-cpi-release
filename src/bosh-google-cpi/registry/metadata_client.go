package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"

	"bosh-google-cpi/google/client"
	opsvc "bosh-google-cpi/google/operation"
	"bosh-google-cpi/util"
)

const (
	GCEMetadataKey       = "bosh_settings"
	metadataClientLogTag = "registryMetadataClient"
	opWaiterRetryMax     = 100
	opMaxSleepExponent   = 3
	opReadyStatus        = "DONE"
)

type metadataGoogleClient interface {
	ComputeService() *compute.Service
	Project() string
}

// MetadataClient represents a GCE metadata client.
type MetadataClient struct {
	options      ClientOptions
	logger       boshlog.Logger
	googleClient metadataGoogleClient
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

	metadata.Items = make([]*compute.MetadataItems, 0, len(i.items))
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
		return AgentSettings{}, bosherr.Errorf("Retrieving settings from instance metadata for instance %q: missing key %q", instanceID, c.options.GCEMetadataKey)
	}
	var agentSettings AgentSettings
	if err = json.Unmarshal([]byte(settings), &agentSettings); err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Unmarshalling agent settings response from instance metadata key %q", c.options.GCEMetadataKey)
	}
	return agentSettings, nil

}

// Update updates the agent settings for a given instance ID. If there are not already agent settings for the instance, it will create ones.
func (c MetadataClient) Update(instanceID string, agentSettings AgentSettings) error {
	settingsJSON, err := json.Marshal(agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling agent settings")
	}
	c.logger.Debug(metadataClientLogTag, "Updating instance metadata for %q", instanceID)

	currentMetadata, err := c.metadata(instanceID)
	if err != nil {
		return err
	}
	currentMetadata.items[c.options.GCEMetadataKey] = string(settingsJSON)

	op, err := c.googleClient.ComputeService().Instances.SetMetadata(c.googleClient.Project(), currentMetadata.zone, instanceID, currentMetadata.computeMetadata()).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating instance metadata with SetMetadata call")
	}
	_, err = c.wait(op)
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating instance metadata with SetMetadata call")
	}
	return nil
}

func (c MetadataClient) metadata(instanceID string) (instanceMetadata, error) {
	svc := c.googleClient.ComputeService().Instances
	listCall := svc.AggregatedList(c.googleClient.Project()).Filter(fmt.Sprintf("name = %q", instanceID))
	var found *compute.Instance
	err := listCall.Pages(context.Background(), func(page *compute.InstanceAggregatedList) error {
		for _, list := range page.Items {
			for _, instance := range list.Instances {
				if instance.Name != instanceID {
					continue
				}
				if found != nil {
					return fmt.Errorf("instance %q is ambiguous across zones", instanceID)
				}
				found = instance
			}
		}
		return nil
	})
	if err != nil {
		return instanceMetadata{}, bosherr.WrapErrorf(err, "Listing instances to find instance %q", instanceID)
	}
	if found == nil {
		return instanceMetadata{}, bosherr.Errorf("Could not find instance %q", instanceID)
	}
	metadata := instanceMetadata{zone: util.ResourceSplitter(found.Zone), items: make(map[string]string)}
	if found.Metadata != nil {
		metadata.fingerprint = found.Metadata.Fingerprint
		for _, item := range found.Metadata.Items {
			if item != nil && item.Value != nil {
				metadata.items[item.Key] = *item.Value
			}
		}
	}
	c.logger.Debug(metadataClientLogTag, "Got metadata for instance %q in zone %q", instanceID, metadata.zone)
	return metadata, nil
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
