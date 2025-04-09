package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"
	"bosh-google-cpi/google/instance"
	instancefakes "bosh-google-cpi/google/instance/fakes"
	registryfakes "bosh-google-cpi/registry/fakes"
)

var _ = Describe("ConfigureNetworks", func() {
	var (
		err      error
		networks Networks

		vmService      *instancefakes.FakeInstanceService
		registryClient *registryfakes.FakeClient

		configureNetworks ConfigureNetworks
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		registryClient = &registryfakes.FakeClient{}
		configureNetworks = NewConfigureNetworks(vmService, registryClient)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			networks = Networks{
				"fake-network-name": &Network{
					Type:    "dynamic",
					IP:      "fake-network-ip",
					Gateway: "fake-network-gateway",
					Netmask: "fake-network-netmask",
					DNS:     []string{"fake-network-dns"},
					Default: []string{"fake-network-default"},
					CloudProperties: NetworkCloudProperties{
						NetworkName:         "fake-network-cloud-network-name",
						Tags:                instance.Tags([]string{"fake-network-cloud-network-tag"}),
						EphemeralExternalIP: true,
						IPForwarding:        false,
					},
				},
			}
		})

		It("returns an error because method is deprecated", func() {
			_, err = configureNetworks.Run("fake-vm-id", networks)
			Expect(err).To(HaveOccurred())
		})
	})
})
