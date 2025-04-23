package action

import (
	"bosh-google-cpi/registry"

	"bosh-google-cpi/google/instance"
)

type Networks map[string]*Network

type Network struct {
	Type            string                 `json:"type,omitempty"`
	IP              string                 `json:"ip,omitempty"`
	Gateway         string                 `json:"gateway,omitempty"`
	Netmask         string                 `json:"netmask,omitempty"`
	DNS             []string               `json:"dns,omitempty"`
	DHCP            bool                   `json:"use_dhcp,omitempty"`
	Default         []string               `json:"default,omitempty"`
	CloudProperties NetworkCloudProperties `json:"cloud_properties,omitempty"`
}

func (ns Networks) AsInstanceServiceNetworks() instance.Networks {
	networks := instance.Networks{}

	for netName, network := range ns {
		networks[netName] = &instance.Network{
			Type:                network.Type,
			IP:                  network.IP,
			Gateway:             network.Gateway,
			Netmask:             network.Netmask,
			DNS:                 network.DNS,
			Default:             network.Default,
			NetworkName:         network.CloudProperties.NetworkName,
			NetworkProjectID:    network.CloudProperties.NetworkProjectID,
			SubnetworkName:      network.CloudProperties.SubnetworkName,
			Tags:                network.CloudProperties.Tags,
			EphemeralExternalIP: network.CloudProperties.EphemeralExternalIP,
			IPForwarding:        network.CloudProperties.IPForwarding,
		}
	}

	return networks
}

func (ns Networks) AsRegistryNetworks() registry.NetworksSettings {
	networksSettings := registry.NetworksSettings{}

	for netName, network := range ns {
		networksSettings[netName] = registry.NetworkSettings{
			Type:    network.Type,
			IP:      network.IP,
			Gateway: network.Gateway,
			Netmask: network.Netmask,
			DHCP:    true,
			DNS:     network.DNS,
			Default: network.Default,
		}
	}

	return networksSettings
}
