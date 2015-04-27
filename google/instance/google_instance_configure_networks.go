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

	if err = i.updateIPForwarding(instance, instanceNetworks); err != nil {
		return err
	}

	if err = i.updateExternalIP(instance, instanceNetworks); err != nil {
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

func (i GoogleInstanceService) updateIPForwarding(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	// If IP Forwarding has changed we need to recreate the VM
	if instance.CanIpForward != instanceNetworks.CanIPForward() {
		i.logger.Debug(googleInstanceServiceLogTag, "Changing IP Forwarding for Google Instance '%s' not supported", instance.Name)
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateExternalIP(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	var err error

	vipNetwork := instanceNetworks.VipNetwork()

	if vipNetwork.IP != "" {
		err = i.updateVipAddress(instance, vipNetwork.IP)
	} else {
		err = i.updateEphemeralExternalIP(instance, instanceNetworks)
	}

	if err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) updateVipAddress(instance *compute.Instance, ipAddress string) error {
	var instanceExternalIP, accessConfigName string
	if len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
		instanceExternalIP = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
		accessConfigName = instance.NetworkInterfaces[0].AccessConfigs[0].Name
	}

	networkInterface := instance.NetworkInterfaces[0].Name

	if instanceExternalIP == "" || instanceExternalIP != ipAddress {
		if instanceExternalIP != "" {
			i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Static IP Address '%s' from Google Instance '%s'", instanceExternalIP, instance.Name)
			err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfigName)
			if err != nil {
				return err
			}
		}

		accessConfig := &compute.AccessConfig{
			Name:  "External NAT",
			Type:  "ONE_TO_ONE_NAT",
			NatIP: ipAddress,
		}

		i.logger.Debug(googleInstanceServiceLogTag, "Attaching Google Static IP Address '%s' to Google Instance '%s'", ipAddress, instance.Name)
		err := i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i GoogleInstanceService) updateEphemeralExternalIP(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	var instanceExternalIP, accessConfigName string
	if len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
		instanceExternalIP = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
		accessConfigName = instance.NetworkInterfaces[0].AccessConfigs[0].Name
	}

	networkInterface := instance.NetworkInterfaces[0].Name

	if instanceNetworks.EphemeralExternalIP() {
		if instanceExternalIP == "" {
			accessConfig := &compute.AccessConfig{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			}

			i.logger.Debug(googleInstanceServiceLogTag, "Attaching Ephemeral Google IP Address to Google Instance '%s'", instance.Name)
			err := i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig)
			if err != nil {
				return err
			}

			return nil
		}

		_, found, err := instanceNetworks.addressService.FindByIP(instanceExternalIP)
		if err != nil {
			return nil
		}

		if found {
			i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Static IP Address '%s' from Google Instance '%s'", instanceExternalIP, instance.Name)
			err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfigName)
			if err != nil {
				return err
			}

			accessConfig := &compute.AccessConfig{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			}

			i.logger.Debug(googleInstanceServiceLogTag, "Attaching Ephemeral Google IP Address to Google Instance '%s'", instance.Name)
			err = i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig)
			if err != nil {
				return err
			}

			return nil
		}
	} else {
		if instanceExternalIP != "" {
			i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Static IP Address '%s' from Google Instance '%s'", instanceExternalIP, instance.Name)
			err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfigName)
			if err != nil {
				return err
			}
		}

		return nil
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
