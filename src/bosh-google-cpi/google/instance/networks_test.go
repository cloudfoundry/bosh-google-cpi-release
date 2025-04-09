package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/google/instance"
)

var _ = Describe("Networks", func() {
	var (
		err            error
		dynamicNetwork *Network
		manualNetwork  *Network
		vipNetwork     *Network
		networks       Networks
	)

	BeforeEach(func() {
		dynamicNetwork = &Network{
			Type:                "dynamic",
			IP:                  "fake-dynamic-network-ip",
			Gateway:             "fake-dynamic-network-gateway",
			Netmask:             "fake-dynamic-network-netmask",
			DNS:                 []string{"fake-dynamic-or-manual-network-dns"},
			Default:             []string{"fake-dynamic-network-default"},
			NetworkName:         "fake-dynamic-network-network-name",
			SubnetworkName:      "fake-dynamic-network-subnetwork-name",
			EphemeralExternalIP: true,
			IPForwarding:        false,
			Tags:                Tags{"fake-dynamic-network-network-tag"},
		}
		manualNetwork = &Network{
			Type:                "manual",
			IP:                  "fake-manual-network-ip",
			Gateway:             "fake-manual-network-gateway",
			Netmask:             "fake-manual-network-netmask",
			DNS:                 []string{"fake-dynamic-or-manual-network-dns"},
			Default:             []string{"fake-manual-network-default"},
			NetworkName:         "fake-manual-network-network-name",
			SubnetworkName:      "fake-manual-network-subnetwork-name",
			EphemeralExternalIP: true,
			IPForwarding:        false,
			Tags:                Tags{"fake-manual-network-network-tag"},
		}
		vipNetwork = &Network{
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
			Tags:                Tags{"fake-vip-network-network-tag"},
		}

		networks = Networks{
			"fake-dynamic-network": dynamicNetwork,
			"fake-vip-network":     vipNetwork,
		}
	})

	Describe("Validate", func() {
		Context("when networks are valid", func() {
			It("does not return an error", func() {
				err = networks.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when networks are not valid", func() {
			BeforeEach(func() {
				networks = Networks{"fake-network-name": &Network{Type: "unknown"}}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when there are not any dynamic or manual networks", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-vip-network": vipNetwork,
				}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Exactly one Dynamic or Manual network must be defined"))
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
				Expect(err.Error()).To(ContainSubstring("Exactly one Dynamic or Manual network must be defined"))
			})
		})

		Context("when there are more than 1 manual networks", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-manual-network-1": manualNetwork,
					"fake-manual-network-2": manualNetwork,
				}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Exactly one Dynamic or Manual network must be defined"))
			})
		})

		Context("when there are more than 1 manual and dynamic networks", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-manual-network-1":  manualNetwork,
					"fake-dynamic-network-2": dynamicNetwork,
				}
			})

			It("returns an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Exactly one Dynamic or Manual network must be defined"))
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
		Context("when there is a dynamic network", func() {
			BeforeEach(func() {
				delete(networks, "fake-manual-network")
			})

			It("returns the dynamic network", func() {
				Expect(networks.Network()).To(Equal(dynamicNetwork))
			})
		})

		Context("when there is NOT a dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("returns an emtpy network", func() {
				Expect(networks.Network()).To(Equal(&Network{}))
			})
		})
	})

	Describe("ManualNetwork", func() {
		Context("when there is a manual network", func() {
			BeforeEach(func() {
				delete(networks, "fake-dynamic-network")
				networks["fake-manual-network"] = manualNetwork
			})
			It("returns the manual network", func() {
				Expect(networks.Network()).To(Equal(manualNetwork))
			})
		})

		Context("when there is NOT a manual network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("returns an emtpy network", func() {
				Expect(networks.Network()).To(Equal(&Network{}))
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
				Expect(networks.VipNetwork()).To(Equal(&Network{}))
			})
		})
	})

	Describe("DNS", func() {
		It("returns only the dynamic network DNS servers", func() {
			Expect(networks.DNS()).To(Equal([]string{"fake-dynamic-or-manual-network-dns"}))
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
			Expect(networks.Tags()).To(Equal(Tags{"fake-dynamic-network-network-tag"}))
		})
	})
})
