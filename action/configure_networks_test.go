package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"

	addressfakes "github.com/frodenas/bosh-google-cpi/google/address_service/fakes"
	instancefakes "github.com/frodenas/bosh-google-cpi/google/instance_service/fakes"
	networkfakes "github.com/frodenas/bosh-google-cpi/google/network_service/fakes"
	targetpoolfakes "github.com/frodenas/bosh-google-cpi/google/target_pool_service/fakes"

	registryfakes "github.com/frodenas/bosh-registry/client/fakes"
)

var _ = Describe("ConfigureNetworks", func() {
	var (
		err               error
		networks          Networks
		vmService         *instancefakes.FakeInstanceService
		addressService    *addressfakes.FakeAddressService
		networkService    *networkfakes.FakeNetworkService
		targetPoolService *targetpoolfakes.FakeTargetPoolService
		registryClient    *registryfakes.FakeClient
		configureNetworks ConfigureNetworks
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		addressService = &addressfakes.FakeAddressService{}
		networkService = &networkfakes.FakeNetworkService{}
		targetPoolService = &targetpoolfakes.FakeTargetPoolService{}
		registryClient = &registryfakes.FakeClient{}
		configureNetworks = NewConfigureNetworks(vmService, addressService, networkService, targetPoolService, registryClient)

		networks = Networks{
			"fake-network-1-name": Network{
				Type:    "dynamic",
				IP:      "fake-network-1-ip",
				Gateway: "fake-network-1-gateway",
				Netmask: "fake-network-1-netmask",
				DNS:     []string{"fake-network-1-dns"},
				Default: []string{"fake-network-1-default"},
				CloudProperties: NetworkCloudProperties{
					NetworkName:         "fake-network-1-cloud-network-name",
					Tags:                NetworkTags{"fake-network-1-cloud-network-tag"},
					EphemeralExternalIP: true,
					IPForwarding:        false,
					TargetPool:          "fake-network-1-cloud-target-pool",
				},
			},
		}
	})

	Describe("Run", func() {
		It("deletes the vm", func() {
			_, err = configureNetworks.Run("fake-vm-id", networks)
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
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
