package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/google/instance_service"
)

var _ = Describe("Network", func() {
	var (
		dynamicNetwork Network
		vipNetwork     Network
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
			EphemeralExternalIP: true,
			IPForwarding:        false,
			Tags:                NetworkTags{"fake-dynamic-network-network-tag"},
			TargetPool:          "fake-dynamic-network-target-pool",
		}

		vipNetwork = Network{
			Type:                "vip",
			IP:                  "fake-vip-network-ip",
			Gateway:             "fake-vip-network-gateway",
			Netmask:             "fake-vip-network-netmask",
			DNS:                 []string{"fake-vip-network-dns"},
			Default:             []string{"fake-vip-network-default"},
			NetworkName:         "fake-vip-network-network-name",
			EphemeralExternalIP: false,
			IPForwarding:        true,
			Tags:                NetworkTags{"fake-vip-network-network-tag"},
			TargetPool:          "fake-vip-network-target-pool",
		}
	})

	Describe("IsDynamic", func() {
		It("returns true for a dynamic network", func() {
			Expect(dynamicNetwork.IsDynamic()).To(BeTrue())
		})

		It("returns false for a vip network", func() {
			Expect(vipNetwork.IsDynamic()).To(BeFalse())
		})
	})

	Describe("IsVip", func() {
		It("returns true for a vip network", func() {
			Expect(vipNetwork.IsVip()).To(BeTrue())
		})

		It("returns false for a dynamic network", func() {
			Expect(dynamicNetwork.IsVip()).To(BeFalse())
		})
	})
})

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
			EphemeralExternalIP: true,
			IPForwarding:        false,
			Tags:                NetworkTags{"fake-dynamic-network-network-tag"},
			TargetPool:          "fake-dynamic-network-target-pool",
		}

		vipNetwork = Network{
			Type:                "vip",
			IP:                  "fake-vip-network-ip",
			Gateway:             "fake-vip-network-gateway",
			Netmask:             "fake-vip-network-netmask",
			DNS:                 []string{"fake-vip-network-dns"},
			Default:             []string{"fake-vip-network-default"},
			NetworkName:         "fake-vip-network-network-name",
			EphemeralExternalIP: false,
			IPForwarding:        true,
			Tags:                NetworkTags{"fake-vip-network-network-tag"},
			TargetPool:          "fake-vip-network-target-pool",
		}

		networks = Networks{
			"fake-dynamic-network": dynamicNetwork,
			"fake-vip-network":     vipNetwork,
		}
	})

	Describe("Validate", func() {
		It("should not return an error", func() {
			err = networks.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is NOT a dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("At least one 'dynamic' network should be defined"))
			})
		})

		Context("when there is more than one dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-dynamic-network-1": dynamicNetwork,
					"fake-dynamic-network-2": dynamicNetwork,
				}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only one dynamic network is allowed"))
			})
		})

		Context("when there is NOT a vip network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-dynamic-network": dynamicNetwork}
			})

			It("should not return an error", func() {
				err = networks.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is more than one vip network", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-vip-network-1": vipNetwork,
					"fake-vip-network-2": vipNetwork,
				}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only one VIP network is allowed"))
			})
		})

		Context("when network tags are not valid", func() {
			BeforeEach(func() {
				dynamicNetwork.Tags = NetworkTags{"invalid_network_tag"}
				networks = Networks{"fake-dynamic-network": dynamicNetwork}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("does not comply with RFC1035"))
			})
		})

		Context("when VIP network does not have an IP", func() {
			BeforeEach(func() {
				vipNetwork.IP = ""
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("VIP Network must have an IP address"))
			})
		})
	})

	Describe("DynamicNetwork", func() {
		It("should return the dynamic network", func() {
			Expect(networks.DynamicNetwork()).To(Equal(dynamicNetwork))
		})

		Context("when there is NOT a dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("should return an emtpy network", func() {
				Expect(networks.DynamicNetwork()).To(Equal(Network{}))
			})
		})
	})

	Describe("VipNetwork", func() {
		It("should return the vip network", func() {
			Expect(networks.VipNetwork()).To(Equal(vipNetwork))
		})

		Context("when there is NOT a vip network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-dynamic-network": dynamicNetwork}
			})

			It("should return an emtpy network", func() {
				Expect(networks.VipNetwork()).To(Equal(Network{}))
			})
		})
	})

	Describe("DNS", func() {
		It("should return only the dynamic network DNS servers", func() {
			Expect(networks.DNS()).To(Equal([]string{"fake-dynamic-network-dns"}))
		})
	})

	Describe("NetworkName", func() {
		It("should return only the dynamic network network name", func() {
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

			It("should return the default network name", func() {
				Expect(networks.NetworkName()).To(Equal("default"))
			})
		})
	})

	Describe("EphemeralExternalIP", func() {
		It("should return only the dynamic network EphemeralExternalIP value", func() {
			Expect(networks.EphemeralExternalIP()).To(BeTrue())
		})
	})

	Describe("CanIPForward", func() {
		It("should return only the dynamic network IPForwarding value", func() {
			Expect(networks.CanIPForward()).To(BeFalse())
		})
	})

	Describe("Tags", func() {
		It("should return only the dynamic network Tags", func() {
			Expect(networks.Tags()).To(Equal(NetworkTags{"fake-dynamic-network-network-tag"}))
		})
	})

	Describe("TargetPool", func() {
		It("should return only the dynamic network TargetPool", func() {
			Expect(networks.TargetPool()).To(Equal("fake-dynamic-network-target-pool"))
		})
	})
})
