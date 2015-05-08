package action

import (
	"github.com/frodenas/bosh-registry/client"

	"github.com/frodenas/bosh-google-cpi/google/instance"
)

type Networks map[string]Network

type Network struct {
	Type            string                 `json:"type,omitempty"`
	IP              string                 `json:"ip,omitempty"`
	Gateway         string                 `json:"gateway,omitempty"`
	Netmask         string                 `json:"netmask,omitempty"`
	DNS             []string               `json:"dns,omitempty"`
	Default         []string               `json:"default,omitempty"`
	CloudProperties NetworkCloudProperties `json:"cloud_properties,omitempty"`
}

func (ns Networks) AsGoogleInstanceNetworks() ginstance.InstanceNetworks {
	networks := ginstance.InstanceNetworks{}

	for netName, network := range ns {
		networks[netName] = ginstance.InstanceNetwork{
			Type:                network.Type,
			IP:                  network.IP,
			Gateway:             network.Gateway,
			Netmask:             network.Netmask,
			DNS:                 network.DNS,
			Default:             network.Default,
			NetworkName:         network.CloudProperties.NetworkName,
			Tags:                ginstance.InstanceNetworkTags(network.CloudProperties.Tags),
			EphemeralExternalIP: network.CloudProperties.EphemeralExternalIP,
			IPForwarding:        network.CloudProperties.IPForwarding,
			TargetPool:          network.CloudProperties.TargetPool,
		}
	}

	return networks
}

func (ns Networks) AsAgentNetworks() registry.NetworksSettings {
	networksSettings := registry.NetworksSettings{}

	for netName, network := range ns {
		networksSettings[netName] = registry.NetworkSettings{
			Type:    network.Type,
			IP:      network.IP,
			Gateway: network.Gateway,
			Netmask: network.Netmask,
			DNS:     network.DNS,
			Default: network.Default,
		}
	}

	return networksSettings
}
