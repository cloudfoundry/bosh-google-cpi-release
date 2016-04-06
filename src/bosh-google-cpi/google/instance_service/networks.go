package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

const defaultNetworkName = "default"

type Networks map[string]Network

func (n Networks) Validate() error {
	var dynamicNetworks, vipNetworks int

	for _, network := range n {
		if err := network.Validate(); err != nil {
			return err
		}

		switch {
		case network.IsDynamic():
			dynamicNetworks++
		case network.IsVip():
			vipNetworks++
		}
	}

	if dynamicNetworks == 0 {
		return bosherr.Error("At least one Dynamic network should be defined")
	}

	if dynamicNetworks > 1 {
		return bosherr.Error("Only one Dynamic network is allowed")
	}

	if vipNetworks > 1 {
		return bosherr.Error("Only one VIP network is allowed")
	}

	return nil
}

func (n Networks) DynamicNetwork() Network {
	for _, net := range n {
		if net.IsDynamic() {
			// There can only be 1 dynamic network
			return net
		}
	}

	return Network{}
}

func (n Networks) VipNetwork() Network {
	for _, net := range n {
		if net.IsVip() {
			// There can only be 1 vip network
			return net
		}
	}

	return Network{}
}

func (n Networks) DNS() []string {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.DNS
}

func (n Networks) NetworkName() string {
	dynamicNetwork := n.DynamicNetwork()

	if dynamicNetwork.NetworkName != "" {
		return dynamicNetwork.NetworkName
	}

	return defaultNetworkName
}

func (n Networks) SubnetworkName() string {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.SubnetworkName
}

func (n Networks) EphemeralExternalIP() bool {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.EphemeralExternalIP
}

func (n Networks) CanIPForward() bool {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.IPForwarding
}

func (n Networks) Tags() NetworkTags {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.Tags
}

func (n Networks) TargetPool() string {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.TargetPool
}

func (n Networks) InstanceGroup() string {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.InstanceGroup
}
