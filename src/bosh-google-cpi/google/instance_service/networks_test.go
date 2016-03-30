package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/google/instance_service"
)

var _ = Describe("Networks", func() {
	var (
		err            error
		dynamicNetwork Network
		vipNetwork     Network
		networks       Networks
	)

	BeforeEach(func() {
		dynamicNetwork = Network{
			Type:                "dynamic",
			IP:                  "fake-dynamic-network-ip",
			Gateway:             "fake-dynamic-network-gateway",
			Netmask:             "fake-dynamic-network-netmask",
			DNS:                 []string{"fake-dynamic-network-dns"},
			Default:             []string{"fake-dynamic-network-default"},
			NetworkName:         "fake-dynamic-network-network-name",
			SubnetworkName:      "fake-dynamic-network-subnetwork-name",
			EphemeralExternalIP: true,
			IPForwarding:        false,
			Tags:                NetworkTags{"fake-dynamic-network-network-tag"},
			TargetPool:          "fake-dynamic-network-target-pool",
			InstanceGroup:       "fake-dynamic-network-instance-group",
		}

		vipNetwork = Network{
			Type:                "vip",
			IP:                  "fake-vip-network-ip",
			Gateway:             "fake-vip-network-gateway",
			Netmask:             "fake-vip-network-netmask",
			DNS:                 []string{"fake-vip-network-dns"},
			Default:             []string{"fake-vip-network-default"},
			NetworkName:         "fake-vip-network-network-name",
			SubnetworkName:      "fake-vip-network-subnetwork-name",
			EphemeralExternalIP: false,
			IPForwarding:        true,
			Tags:                NetworkTags{"fake-vip-network-network-tag"},
			TargetPool:          "fake-vip-network-target-pool",
			InstanceGroup:       "fake-vip-network-instance-group",
		}

		networks = Networks{
			"fake-dynamic-network": dynamicNetwork,
			"fake-vip-network":     vipNetwork,
		}
	})

	Describe("Validate", func() {
		It("does not return an error if networks are valid", func() {
			err = networks.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when networks are not valid", func() {
			BeforeEach(func() {
				networks = Networks{"fake-network-name": Network{Type: "unknown"}}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when there are not any dynamic networks", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-vip-network": vipNetwork,
				}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("At least one Dynamic network should be defined"))
			})
		})

		Context("when there are more than 1 dynamic networks", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-dynamic-network-1": dynamicNetwork,
					"fake-dynamic-network-2": dynamicNetwork,
				}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only one Dynamic network is allowed"))
			})
		})

		Context("when there are more than 1 vip networks", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-dynamic-network": dynamicNetwork,
					"fake-vip-network-1":   vipNetwork,
					"fake-vip-network-2":   vipNetwork,
				}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only one VIP network is allowed"))
			})
		})
	})

	Describe("DynamicNetwork", func() {
		It("returns the dynamic network", func() {
			Expect(networks.DynamicNetwork()).To(Equal(dynamicNetwork))
		})

		Context("when there is NOT a dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("returns an emtpy network", func() {
				Expect(networks.DynamicNetwork()).To(Equal(Network{}))
			})
		})
	})

	Describe("VipNetwork", func() {
		It("returns the vip network", func() {
			Expect(networks.VipNetwork()).To(Equal(vipNetwork))
		})

		Context("when there is NOT a vip network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-dynamic-network": dynamicNetwork}
			})

			It("returns an emtpy network", func() {
				Expect(networks.VipNetwork()).To(Equal(Network{}))
			})
		})
	})

	Describe("DNS", func() {
		It("returns only the dynamic network DNS servers", func() {
			Expect(networks.DNS()).To(Equal([]string{"fake-dynamic-network-dns"}))
		})
	})

	Describe("NetworkName", func() {
		It("returns only the dynamic network network name", func() {
			Expect(networks.NetworkName()).To(Equal("fake-dynamic-network-network-name"))
		})

		Context("when Network Name is empty", func() {
			BeforeEach(func() {
				dynamicNetwork.NetworkName = ""
				networks = Networks{
					"fake-dynamic-network": dynamicNetwork,
					"fake-vip-network":     vipNetwork,
				}
			})

			It("returns the default network name", func() {
				Expect(networks.NetworkName()).To(Equal("default"))
			})
		})
	})

	Describe("SubnetworkName", func() {
		It("returns only the dynamic network subnetwork name", func() {
			Expect(networks.SubnetworkName()).To(Equal("fake-dynamic-network-subnetwork-name"))
		})

		Context("when Subnetwork Name is empty", func() {
			BeforeEach(func() {
				dynamicNetwork.SubnetworkName = ""
				networks = Networks{
					"fake-dynamic-network": dynamicNetwork,
					"fake-vip-network":     vipNetwork,
				}
			})

			It("returns an empty subnetwork name", func() {
				Expect(networks.SubnetworkName()).To(BeEmpty())
			})
		})
	})

	Describe("EphemeralExternalIP", func() {
		It("returns only the dynamic network EphemeralExternalIP value", func() {
			Expect(networks.EphemeralExternalIP()).To(BeTrue())
		})
	})

	Describe("CanIPForward", func() {
		It("returns only the dynamic network IPForwarding value", func() {
			Expect(networks.CanIPForward()).To(BeFalse())
		})
	})

	Describe("Tags", func() {
		It("returns only the dynamic network Tags", func() {
			Expect(networks.Tags()).To(Equal(NetworkTags{"fake-dynamic-network-network-tag"}))
		})
	})

	Describe("TargetPool", func() {
		It("returns only the dynamic network TargetPool", func() {
			Expect(networks.TargetPool()).To(Equal("fake-dynamic-network-target-pool"))
		})
	})

	Describe("InstanceGroup", func() {
		It("returns only the dynamic network InstanceGroup", func() {
			Expect(networks.InstanceGroup()).To(Equal("fake-dynamic-network-instance-group"))
		})
	})
})
