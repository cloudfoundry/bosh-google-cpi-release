package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/google/instance_service"
)

var _ = Describe("Network", func() {
	var (
		err            error
		dynamicNetwork Network
		vipNetwork     Network
		unknownNetwork Network
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

		unknownNetwork = Network{Type: "unknown"}
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

	Describe("Validate", func() {
		Context("Dynamic Network", func() {
			It("does not return error if network properties are valid", func() {
				err = dynamicNetwork.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when network tags are not valid", func() {
				BeforeEach(func() {
					dynamicNetwork.Tags = NetworkTags{"invalid_network_tag"}
				})

				It("returns an error", func() {
					err = dynamicNetwork.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("does not comply with RFC1035"))
				})
			})
		})

		Context("VIP Network", func() {
			It("does not return error if network properties are valid", func() {
				err = vipNetwork.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error if does not have n IP Address", func() {
				vipNetwork.IP = ""

				err = vipNetwork.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("VIP Networks must provide an IP Address"))
			})
		})

		Context("Unknown Network", func() {
			It("returns an error", func() {
				err = unknownNetwork.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Network type 'unknown' not supported"))
			})
		})
	})
})
