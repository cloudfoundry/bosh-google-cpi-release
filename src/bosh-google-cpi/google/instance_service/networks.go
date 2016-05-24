package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

const defaultNetworkName = "default"

type Networks map[string]Network

func (n Networks) Validate() error {
	var networks, vipNetworks int

	for _, network := range n {
		if err := network.Validate(); err != nil {
			return err
		}

		switch {
		case network.IsDynamic():
			networks++
		case network.IsManual():
			networks++
		case network.IsVip():
			vipNetworks++
		}
	}

	if networks != 1 {
		return bosherr.Error("Exactly one Dynamic or Manual network must be defined")
	}

	if vipNetworks > 1 {
		return bosherr.Error("Only one VIP network is allowed")
	}

	return nil
}

func (n Networks) Network() Network {
	for _, net := range n {
		if !net.IsVip() {
			// There can only be 1 dynamic or manual network
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
	network := n.Network()

	return network.DNS
}

func (n Networks) NetworkName() string {
	network := n.Network()

	if network.NetworkName != "" {
		return network.NetworkName
	}

	return defaultNetworkName
}

func (n Networks) SubnetworkName() string {
	network := n.Network()

	return network.SubnetworkName
}

func (n Networks) EphemeralExternalIP() bool {
	network := n.Network()

	return network.EphemeralExternalIP
}

func (n Networks) StaticPrivateIP() string {
	network := n.Network()

	return network.IP
}

func (n Networks) CanIPForward() bool {
	network := n.Network()

	return network.IPForwarding
}

func (n Networks) Tags() NetworkTags {
	network := n.Network()

	return network.Tags
}
