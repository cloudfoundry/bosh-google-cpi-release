package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	instancefakes "bosh-google-cpi/google/instance_service/fakes"

	registryfakes "github.com/frodenas/bosh-registry/client/fakes"

	"github.com/frodenas/bosh-registry/client"
)

var _ = Describe("ConfigureNetworks", func() {
	var (
		err                   error
		networks              Networks
		expectedAgentSettings registry.AgentSettings

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
				"fake-network-name": Network{
					Type:    "dynamic",
					IP:      "fake-network-ip",
					Gateway: "fake-network-gateway",
					Netmask: "fake-network-netmask",
					DNS:     []string{"fake-network-dns"},
					Default: []string{"fake-network-default"},
					CloudProperties: NetworkCloudProperties{
						NetworkName:         "fake-network-cloud-network-name",
						Tags:                NetworkTags{"fake-network-cloud-network-tag"},
						EphemeralExternalIP: true,
						IPForwarding:        false,
						TargetPool:          "fake-network-cloud-target-pool",
					},
				},
			}
			registryClient.FetchSettings = registry.AgentSettings{}
			expectedAgentSettings = registry.AgentSettings{
				Networks: registry.NetworksSettings{
					"fake-network-name": registry.NetworkSettings{
						Type:    "dynamic",
						IP:      "fake-network-ip",
						Gateway: "fake-network-gateway",
						Netmask: "fake-network-netmask",
						DNS:     []string{"fake-network-dns"},
						Default: []string{"fake-network-default"},
					},
				},
			}
		})

		It("configures the network", func() {
			_, err = configureNetworks.Run("fake-vm-id", networks)
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
		})

		It("returns an error if networs are not valied", func() {
			_, err = configureNetworks.Run("fake-vm-id", Networks{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Configuring networks for vm"))
			Expect(vmService.UpdateNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if vmService update network configuration call returns an error", func() {
			vmService.UpdateNetworkConfigurationErr = errors.New("fake-vm-service-error")

			_, err = configureNetworks.Run("fake-vm-id", networks)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient fetch call returns an error", func() {
			registryClient.FetchErr = errors.New("fake-registry-client-error")

			_, err = configureNetworks.Run("fake-vm-id", networks)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(vmService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = configureNetworks.Run("fake-vm-id", networks)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(vmService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})
	})
})
