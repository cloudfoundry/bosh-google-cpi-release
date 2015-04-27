package ginstance

import (
	"reflect"
	"sort"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceService) ConfigureNetworks(id string, instanceNetworks GoogleInstanceNetworks) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Instance '%s' not found", id)
	}

	// TODO: Configure VIP network

	if err := i.addToTargetPool(instance, instanceNetworks); err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) addToTargetPool(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	targetPoolName := instanceNetworks.TargetPool()

	if targetPoolName != "" {
		err := instanceNetworks.targetPoolService.AddInstance(targetPoolName, instance.SelfLink)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i GoogleInstanceService) UpdateNetworks(id string, instanceNetworks GoogleInstanceNetworks) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Instance '%s' not found", id)
	}

	if err = i.updateNetwork(instance, instanceNetworks); err != nil {
		return err
	}

	if err = i.updateIpForwarding(instance, instanceNetworks); err != nil {
		return err
	}

	if err = i.updateEphemeralExternalIp(instance, instanceNetworks); err != nil {
		return err
	}

	if err = i.updateTags(instance, instanceNetworks); err != nil {
		return err
	}

	if err := i.updateTargetPool(instance, instanceNetworks); err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) updateNetwork(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	// If the network has changed we need to recreate the VM
	dynamicNetwork := instanceNetworks.DynamicNetwork()
	if gutil.ResourceSplitter(instance.NetworkInterfaces[0].Network) != dynamicNetwork.NetworkName {
		i.logger.Debug(googleInstanceServiceLogTag, "Changing network for Google Instance '%s' not supported", instance.Name)
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateIpForwarding(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	// If IP Forwarding has changed we need to recreate the VM
	if instance.CanIpForward != instanceNetworks.CanIpForward() {
		i.logger.Debug(googleInstanceServiceLogTag, "Changing IP Forwarding for Google Instance '%s' not supported", instance.Name)
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateEphemeralExternalIp(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	var instanceExternalIp string

	if len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
		instanceExternalIp = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
	}

	if instanceNetworks.EphemeralExternalIP() {
		if instanceExternalIp == "" {
			networkInterface := instance.NetworkInterfaces[0].Name
			accessConfig := &compute.AccessConfig{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			}
			err := i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig)
			if err != nil {
				return err
			}
		}
	} else {
		if instanceExternalIp != "" {
			// TODO: Only if network has no vip
			networkInterface := instance.NetworkInterfaces[0].Name
			accessConfig := instance.NetworkInterfaces[0].AccessConfigs[0].Name
			err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i GoogleInstanceService) updateTags(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	// Parset network tags
	networkTags, err := instanceNetworks.Tags()
	if err != nil {
		return err
	}

	// Check if tags have changed
	sort.Strings(networkTags.Items)
	sort.Strings(instance.Tags.Items)
	if reflect.DeepEqual(networkTags.Items, instance.Tags.Items) {
		return nil
	}

	// Override the instance tags preserving the original fingerprint
	instanceTags := &compute.Tags{
		Fingerprint: instance.Tags.Fingerprint,
		Items:       networkTags.Items,
	}

	// Update the instance tags
	err = i.SetTags(instance.Name, instance.Zone, instanceTags)
	if err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) updateTargetPool(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	targetPoolName := instanceNetworks.TargetPool()
	currentTargetPool, _, err := instanceNetworks.targetPoolService.FindByInstance(instance.SelfLink, "")
	if err != nil {
		return err
	}

	if targetPoolName != currentTargetPool {
		if currentTargetPool != "" {
			err := instanceNetworks.targetPoolService.RemoveInstance(currentTargetPool, instance.SelfLink)
			if err != nil {
				return err
			}
		}

		if targetPoolName != "" {
			err := instanceNetworks.targetPoolService.AddInstance(targetPoolName, instance.SelfLink)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
