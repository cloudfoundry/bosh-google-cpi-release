package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	instancefakes "bosh-google-cpi/google/instance_service/fakes"

	registryfakes "github.com/frodenas/bosh-registry/client/fakes"
)

var _ = Describe("DeleteVM", func() {
	var (
		err error

		vmService      *instancefakes.FakeInstanceService
		registryClient *registryfakes.FakeClient

		deleteVM DeleteVM
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		registryClient = &registryfakes.FakeClient{}
		deleteVM = NewDeleteVM(vmService, registryClient)
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
