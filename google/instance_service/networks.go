package instance

import (
	"regexp"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

const defaultNetworkName = "default"
const maxTagLength = 63

type Networks map[string]Network

type Network struct {
	Type                string
	IP                  string
	Gateway             string
	Netmask             string
	DNS                 []string
	Default             []string
	NetworkName         string
	EphemeralExternalIP bool
	IPForwarding        bool
	Tags                NetworkTags
	TargetPool          string
}

type NetworkTags []string

func (n Network) IsDynamic() bool { return n.Type == "dynamic" }

func (n Network) validateDynamic() error {
	err := n.validateTags()
	if err != nil {
		return err
	}

	return nil
}

func (n Network) validateTags() error {
	if len(n.Tags) > 0 {
		pattern, _ := regexp.Compile("^[A-Za-z]+[A-Za-z0-9-]*[A-Za-z0-9]+$")
		for _, tag := range n.Tags {
			if len(tag) > maxTagLength || !pattern.MatchString(tag) {
				return bosherr.Errorf("Invalid tag '%s': does not comply with RFC1035", tag)
			}
		}
	}

	return nil
}

func (n Network) IsVip() bool { return n.Type == "vip" }

func (n Network) validateVip() error {
	if n.IP == "" {
		return bosherr.Error("VIP Network must have an IP address")
	}

	return nil
}

func (n Networks) Validate() error {
	var dnet, vnet bool

	for _, net := range n {
		if net.IsDynamic() {
			if dnet {
				return bosherr.Error("Only one dynamic network is allowed")
			}

			err := net.validateDynamic()
			if err != nil {
				return err
			}

			dnet = true
		}

		if net.IsVip() {
			if vnet {
				return bosherr.Error("Only one VIP network is allowed")
			}

			err := net.validateVip()
			if err != nil {
				return err
			}

			vnet = true
		}
	}

	if !dnet {
		return bosherr.Error("At least one 'dynamic' network should be defined")
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
