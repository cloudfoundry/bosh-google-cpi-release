package ginstance

import (
	"regexp"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/address_service"
	"github.com/frodenas/bosh-google-cpi/google/network_service"
	"github.com/frodenas/bosh-google-cpi/google/target_pool_service"
	"google.golang.org/api/compute/v1"
)

const googleMaxTagLength = 63

type GoogleInstanceNetworks struct {
	networks          InstanceNetworks
	addressService    gaddress.AddressService
	networkService    gnetwork.NetworkService
	targetPoolService gtargetpool.TargetPoolService
}

func NewGoogleInstanceNetworks(
	networks InstanceNetworks,
	addressService gaddress.AddressService,
	networkService gnetwork.NetworkService,
	targetPoolService gtargetpool.TargetPoolService,
) GoogleInstanceNetworks {
	return GoogleInstanceNetworks{
		networks:          networks,
		addressService:    addressService,
		networkService:    networkService,
		targetPoolService: targetPoolService,
	}
}

func (in GoogleInstanceNetworks) DynamicNetwork() InstanceNetwork {
	for _, net := range in.networks {
		if net.IsDynamic() {
			// There can only be 1 dynamic network
			return net
		}
	}

	return InstanceNetwork{}
}

func (in GoogleInstanceNetworks) VipNetwork() InstanceNetwork {
	for _, net := range in.networks {
		if net.IsVip() {
			// There can only be 1 vip network
			return net
		}
	}

	return InstanceNetwork{}
}

func (in GoogleInstanceNetworks) CanIPForward() bool {
	dynamicNetwork := in.DynamicNetwork()

	return dynamicNetwork.IPForwarding
}

func (in GoogleInstanceNetworks) DNS() []string {
	dynamicNetwork := in.DynamicNetwork()

	return dynamicNetwork.DNS
}

func (in GoogleInstanceNetworks) EphemeralExternalIP() bool {
	dynamicNetwork := in.DynamicNetwork()

	return dynamicNetwork.EphemeralExternalIP
}

func (in GoogleInstanceNetworks) NetworkInterfaces() ([]*compute.NetworkInterface, error) {
	dynamicNetwork := in.DynamicNetwork()
	network, found, err := in.networkService.Find(dynamicNetwork.NetworkName)
	if err != nil {
		return nil, bosherr.WrapError(err, "Network Interfaces")
	}
	if !found {
		return nil, bosherr.WrapErrorf(err, "Network Interfaces: Network '%s' does not exists", dynamicNetwork.NetworkName)
	}

	var networkInterfaces []*compute.NetworkInterface
	var accessConfigs []*compute.AccessConfig

	vipNetwork := in.VipNetwork()
	if dynamicNetwork.EphemeralExternalIP || vipNetwork.IP != "" {
		accessConfig := &compute.AccessConfig{
			Name: "External NAT",
			Type: "ONE_TO_ONE_NAT",
		}
		if vipNetwork.IP != "" {
			accessConfig.NatIP = vipNetwork.IP
		}
		accessConfigs = append(accessConfigs, accessConfig)
	}

	networkInterface := &compute.NetworkInterface{
		Network:       network.SelfLink,
		AccessConfigs: accessConfigs,
	}
	networkInterfaces = append(networkInterfaces, networkInterface)

	return networkInterfaces, nil
}

func (in GoogleInstanceNetworks) Tags() (*compute.Tags, error) {
	tags := &compute.Tags{}

	dynamicNetwork := in.DynamicNetwork()
	pattern, _ := regexp.Compile("^[A-Za-z]+[A-Za-z0-9-]*[A-Za-z0-9]+$")
	for _, tag := range dynamicNetwork.Tags {
		if len(tag) > googleMaxTagLength || !pattern.MatchString(tag) {
			return tags, bosherr.Errorf("Invalid tag '%s': does not comply with RFC1035", tag)
		}
		tags.Items = append(tags.Items, tag)
	}

	return tags, nil
}

func (in GoogleInstanceNetworks) TargetPool() string {
	dynamicNetwork := in.DynamicNetwork()

	return dynamicNetwork.TargetPool
}

func (in GoogleInstanceNetworks) Validate() error {
	var dnet, vnet bool

	// TODO: refactor & add more validations
	for _, net := range in.networks {
		if net.IsDynamic() {
			if dnet {
				return bosherr.Error("Only one dynamic network is allowed")
			}
			dnet = true
		}

		if net.IsVip() {
			if vnet {
				return bosherr.Error("Only one VIP network is allowed")
			}
			vnet = true
		}
	}

	if !dnet {
		return bosherr.Error("At least one 'dynamic' network should be defined")
	}

	return nil
}
