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

var _ = Describe("DeleteVM", func() {
	var (
		err error

		vmService         *instancefakes.FakeInstanceService
		addressService    *addressfakes.FakeAddressService
		networkService    *networkfakes.FakeNetworkService
		targetPoolService *targetpoolfakes.FakeTargetPoolService
		registryClient    *registryfakes.FakeClient

		deleteVM DeleteVM
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		addressService = &addressfakes.FakeAddressService{}
		networkService = &networkfakes.FakeNetworkService{}
		targetPoolService = &targetpoolfakes.FakeTargetPoolService{}
		registryClient = &registryfakes.FakeClient{}
		deleteVM = NewDeleteVM(vmService, addressService, networkService, targetPoolService, registryClient)
	})

	Describe("Run", func() {
		It("deletes the vm", func() {
			_, err = deleteVM.Run("fake-vm-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(vmService.DeleteCalled).To(BeTrue())
			Expect(registryClient.DeleteCalled).To(BeTrue())
		})

		It("returns an error if vmService delete network configuration call returns an error", func() {
			vmService.DeleteNetworkConfigurationErr = errors.New("fake-vm-service-error")

			_, err = deleteVM.Run("fake-vm-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(vmService.DeleteCalled).To(BeFalse())
			Expect(registryClient.DeleteCalled).To(BeFalse())
		})

		It("returns an error if vmService delete call returns an error", func() {
			vmService.DeleteErr = errors.New("fake-vm-service-error")

			_, err = deleteVM.Run("fake-vm-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(vmService.DeleteCalled).To(BeTrue())
			Expect(registryClient.DeleteCalled).To(BeFalse())
		})

		It("returns an error if registryClient delete call returns an error", func() {
			registryClient.DeleteErr = errors.New("fake-registry-client-error")

			_, err = deleteVM.Run("fake-vm-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(vmService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(vmService.DeleteCalled).To(BeTrue())
			Expect(registryClient.DeleteCalled).To(BeTrue())
		})
	})
})
