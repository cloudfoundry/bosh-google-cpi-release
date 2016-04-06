package instance

import (
	"reflect"
	"sort"

	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceService) AddNetworkConfiguration(id string, networks Networks) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	if err := i.addToTargetPool(instance, networks); err != nil {
		return err
	}

	if err := i.addToInstanceGroup(instance, networks); err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) addToTargetPool(instance *compute.Instance, networks Networks) error {
	if targetPoolName := networks.TargetPool(); targetPoolName != "" {
		if err := i.targetPoolService.AddInstance(targetPoolName, instance.SelfLink); err != nil {
			return err
		}
	}

	return nil
}

func (i GoogleInstanceService) addToInstanceGroup(instance *compute.Instance, networks Networks) error {
	if instanceGroupName := networks.InstanceGroup(); instanceGroupName != "" {
		if err := i.instanceGroupService.AddInstance(instanceGroupName, instance.SelfLink); err != nil {
			return err
		}
	}

	return nil
}

func (i GoogleInstanceService) DeleteNetworkConfiguration(id string) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	if err := i.removeFromTargetPool(instance); err != nil {
		return err
	}

	if err := i.removeFromInstanceGroup(instance); err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) removeFromTargetPool(instance *compute.Instance) error {
	targetPool, found, err := i.targetPoolService.FindByInstance(instance.SelfLink, "")
	if err != nil {
		return err
	}

	if found {
		if err := i.targetPoolService.RemoveInstance(targetPool, instance.SelfLink); err != nil {
			return err
		}
	}

	return nil
}

func (i GoogleInstanceService) removeFromInstanceGroup(instance *compute.Instance) error {
	instanceGroup, found, err := i.instanceGroupService.FindByInstance(instance.SelfLink, "")
	if err != nil {
		return err
	}

	if found {
		if err := i.instanceGroupService.RemoveInstance(instanceGroup, instance.SelfLink); err != nil {
			return err
		}
	}

	return nil
}

func (i GoogleInstanceService) UpdateNetworkConfiguration(id string, networks Networks) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	if err = i.updateNetwork(instance, networks); err != nil {
		return err
	}

	if err = i.updateSubnetwork(instance, networks); err != nil {
		return err
	}

	if err = i.updateIPForwarding(instance, networks); err != nil {
		return err
	}

	if err = i.updateExternalIP(instance, networks); err != nil {
		return err
	}

	if err = i.updateTags(instance, networks); err != nil {
		return err
	}

	if err := i.updateTargetPool(instance, networks); err != nil {
		return err
	}

	if err := i.updateInstanceGroup(instance, networks); err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) updateNetwork(instance *compute.Instance, networks Networks) error {
	// If the network has changed we need to recreate the VM
	if util.ResourceSplitter(instance.NetworkInterfaces[0].Network) != networks.NetworkName() {
		i.logger.Debug(googleInstanceServiceLogTag, "Changing network for Google Instance '%s' not supported", instance.Name)
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateSubnetwork(instance *compute.Instance, networks Networks) error {
	if networks.SubnetworkName() == "" {
		return nil
	}

	// If the subnetwork has changed we need to recreate the VM
	if util.ResourceSplitter(instance.NetworkInterfaces[0].Subnetwork) != networks.SubnetworkName() {
		i.logger.Debug(googleInstanceServiceLogTag, "Changing subnetwork for Google Instance '%s' not supported", instance.Name)
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateIPForwarding(instance *compute.Instance, networks Networks) error {
	// If IP Forwarding has changed we need to recreate the VM
	if instance.CanIpForward != networks.CanIPForward() {
		i.logger.Debug(googleInstanceServiceLogTag, "Changing IP Forwarding for Google Instance '%s' not supported", instance.Name)
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateExternalIP(instance *compute.Instance, networks Networks) error {
	var err error

	vipNetwork := networks.VipNetwork()
	if vipNetwork.IP != "" {
		err = i.updateVipAddress(instance, vipNetwork.IP)
	} else {
		err = i.updateEphemeralExternalIP(instance, networks)
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
		// If the instance has an external IP, detach it
		if instanceExternalIP != "" {
			i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Static IP Address '%s' from Google Instance '%s'", instanceExternalIP, instance.Name)
			if err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfigName); err != nil {
				return err
			}
		}

		// Attach the vip IP to the instance
		accessConfig := &compute.AccessConfig{
			Name:  "External NAT",
			Type:  "ONE_TO_ONE_NAT",
			NatIP: ipAddress,
		}

		i.logger.Debug(googleInstanceServiceLogTag, "Attaching Google Static IP Address '%s' to Google Instance '%s'", ipAddress, instance.Name)
		if err := i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig); err != nil {
			return err
		}
	}

	return nil
}

func (i GoogleInstanceService) updateEphemeralExternalIP(instance *compute.Instance, networks Networks) error {
	var instanceExternalIP, accessConfigName string
	if len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
		instanceExternalIP = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
		accessConfigName = instance.NetworkInterfaces[0].AccessConfigs[0].Name
	}

	networkInterface := instance.NetworkInterfaces[0].Name

	if networks.EphemeralExternalIP() {
		// If the instance doesn't have an external IP, attach an ephemeral one
		if instanceExternalIP == "" {
			accessConfig := &compute.AccessConfig{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			}

			i.logger.Debug(googleInstanceServiceLogTag, "Attaching Ephemeral Google IP Address to Google Instance '%s'", instance.Name)
			if err := i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig); err != nil {
				return err
			}

			return nil
		}

		// Check if the instance external IP is an static IP address
		_, found, err := i.addressService.FindByIP(instanceExternalIP)
		if err != nil {
			return nil
		}

		if found {
			// Detach the static IP from the instance
			i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Static IP Address '%s' from Google Instance '%s'", instanceExternalIP, instance.Name)
			if err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfigName); err != nil {
				return err
			}

			// Attach an ephemeral IP to the instance
			accessConfig := &compute.AccessConfig{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			}

			i.logger.Debug(googleInstanceServiceLogTag, "Attaching Ephemeral Google IP Address to Google Instance '%s'", instance.Name)
			if err = i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig); err != nil {
				return err
			}

			return nil
		}
	} else {
		// If the instance has an external IP, detach it from the instance
		if instanceExternalIP != "" {
			i.logger.Debug(googleInstanceServiceLogTag, "Detaching Google Static IP Address '%s' from Google Instance '%s'", instanceExternalIP, instance.Name)
			if err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfigName); err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func (i GoogleInstanceService) updateTags(instance *compute.Instance, networks Networks) error {
	// Parset network tags
	networkTags := networks.Tags()

	// Check if tags have changed
	sort.Strings(networkTags)
	sort.Strings(instance.Tags.Items)
	if reflect.DeepEqual(networkTags, instance.Tags.Items) {
		return nil
	}

	// Override the instance tags preserving the original fingerprint
	instanceTags := &compute.Tags{
		Fingerprint: instance.Tags.Fingerprint,
		Items:       networkTags,
	}

	// Update the instance tags
	if err := i.SetTags(instance.Name, instance.Zone, instanceTags); err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) updateTargetPool(instance *compute.Instance, networks Networks) error {
	// Check if instance is associated to a target pool
	currentTargetPool, _, err := i.targetPoolService.FindByInstance(instance.SelfLink, "")
	if err != nil {
		return err
	}

	// Check if target pool info has changed
	targetPoolName := networks.TargetPool()
	if targetPoolName != currentTargetPool {
		if currentTargetPool != "" {
			if err := i.targetPoolService.RemoveInstance(currentTargetPool, instance.SelfLink); err != nil {
				return err
			}
		}

		if targetPoolName != "" {
			if err := i.targetPoolService.AddInstance(targetPoolName, instance.SelfLink); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i GoogleInstanceService) updateInstanceGroup(instance *compute.Instance, networks Networks) error {
	// Check if instance is associated to an instance group
	currentInstanceGroup, _, err := i.instanceGroupService.FindByInstance(instance.SelfLink, "")
	if err != nil {
		return err
	}

	// Check if instance group info has changed
	instanceGroupName := networks.InstanceGroup()
	if instanceGroupName != currentInstanceGroup {
		if currentInstanceGroup != "" {
			if err := i.instanceGroupService.RemoveInstance(currentInstanceGroup, instance.SelfLink); err != nil {
				return err
			}
		}

		if instanceGroupName != "" {
			if err := i.instanceGroupService.AddInstance(instanceGroupName, instance.SelfLink); err != nil {
				return err
			}
		}
	}

	return nil
}
