package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	instancefakes "bosh-google-cpi/google/instance_service/fakes"
)

var _ = Describe("RebootVM", func() {
	var (
		err error

		vmService *instancefakes.FakeInstanceService

		rebootVM RebootVM
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		rebootVM = NewRebootVM(vmService)
	})

	Describe("Run", func() {
		It("reboots the vm", func() {
			_, err = rebootVM.Run("fake-vm-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.RebootCalled).To(BeTrue())
		})

		It("returns an error if vmService reboot call returns an error", func() {
			vmService.RebootErr = errors.New("fake-vm-service-error")

			_, err = rebootVM.Run("fake-vm-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.RebootCalled).To(BeTrue())
		})
	})
})
