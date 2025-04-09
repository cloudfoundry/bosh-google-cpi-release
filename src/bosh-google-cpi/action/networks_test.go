package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"
	"bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/registry"
)

var _ = Describe("Networks", func() {
	var (
		networks Networks
	)

	BeforeEach(func() {
		networks = Networks{
			"fake-network-1-name": &Network{
				Type:    "fake-network-1-type",
				IP:      "fake-network-1-ip",
				Gateway: "fake-network-1-gateway",
				Netmask: "fake-network-1-netmask",
				DNS:     []string{"fake-network-1-dns"},
				DHCP:    true,
				Default: []string{"fake-network-1-default"},
				CloudProperties: NetworkCloudProperties{
					NetworkName:         "fake-network-1-cloud-network-name",
					NetworkProjectID:    "fake-network-1-cloud-network-project-id",
					SubnetworkName:      "fake-network-1-cloud-subnetwork-name",
					Tags:                instance.Tags([]string{"fake-network-1-cloud-network-tag"}),
					EphemeralExternalIP: true,
					IPForwarding:        false,
				},
			},
			"fake-network-2-name": &Network{
				Type: "fake-network-2-type",
				IP:   "fake-network-2-ip",
				DHCP: true,
			},
		}
	})

	Describe("AsInstanceServiceNetworks", func() {
		It("returns networks for the instance service", func() {
			expectedInstanceNetworks := instance.Networks{
				"fake-network-1-name": &instance.Network{
					Type:                "fake-network-1-type",
					IP:                  "fake-network-1-ip",
					Gateway:             "fake-network-1-gateway",
					Netmask:             "fake-network-1-netmask",
					DNS:                 []string{"fake-network-1-dns"},
					Default:             []string{"fake-network-1-default"},
					NetworkName:         "fake-network-1-cloud-network-name",
					NetworkProjectID:    "fake-network-1-cloud-network-project-id",
					SubnetworkName:      "fake-network-1-cloud-subnetwork-name",
					Tags:                instance.Tags([]string{"fake-network-1-cloud-network-tag"}),
					EphemeralExternalIP: true,
					IPForwarding:        false,
				},
				"fake-network-2-name": &instance.Network{
					Type: "fake-network-2-type",
					IP:   "fake-network-2-ip",
				},
			}

			Expect(networks.AsInstanceServiceNetworks()).To(Equal(expectedInstanceNetworks))
		})
	})

	Describe("AsRegistryNetworks", func() {
		It("returns networks for the registry", func() {
			expectedRegistryNetworks := registry.NetworksSettings{
				"fake-network-1-name": registry.NetworkSettings{
					Type:    "fake-network-1-type",
					IP:      "fake-network-1-ip",
					Gateway: "fake-network-1-gateway",
					Netmask: "fake-network-1-netmask",
					DNS:     []string{"fake-network-1-dns"},
					DHCP:    true,
					Default: []string{"fake-network-1-default"},
				},
				"fake-network-2-name": registry.NetworkSettings{
					Type: "fake-network-2-type",
					IP:   "fake-network-2-ip",
					DHCP: true,
				},
			}

			Expect(networks.AsRegistryNetworks()).To(Equal(expectedRegistryNetworks))
		})
	})
})
