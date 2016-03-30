package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	instancefakes "bosh-google-cpi/google/instance_service/fakes"
)

var _ = Describe("HasVM", func() {
	var (
		err   error
		found bool

		vmService *instancefakes.FakeInstanceService

		hasVM HasVM
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		hasVM = NewHasVM(vmService)
	})

	Describe("Run", func() {
		It("returns true if vm ID exist", func() {
			vmService.FindFound = true

			found, err = hasVM.Run("fake-vm-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(vmService.FindCalled).To(BeTrue())
		})

		It("returns false if vm ID does not exist", func() {
			vmService.FindFound = false

			found, err = hasVM.Run("fake-vm-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(vmService.FindCalled).To(BeTrue())
		})

		It("returns an error if vmService find call returns an error", func() {
			vmService.FindErr = errors.New("fake-vm-service-error")

			_, err = hasVM.Run("fake-vm-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.FindCalled).To(BeTrue())
		})
	})
})
