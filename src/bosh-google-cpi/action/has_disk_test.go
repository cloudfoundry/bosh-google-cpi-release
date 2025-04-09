package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"
	diskfakes "bosh-google-cpi/google/disk/fakes"
)

var _ = Describe("HasDisk", func() {
	var (
		err   error
		found bool

		diskService *diskfakes.FakeDiskService

		hasDisk HasDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		hasDisk = NewHasDisk(diskService)
	})

	Describe("Run", func() {
		It("returns true if disk ID exist", func() {
			diskService.FindFound = true

			found, err = hasDisk.Run("fake-disk-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(diskService.FindCalled).To(BeTrue())
		})

		It("returns false if disk ID does not exist", func() {
			diskService.FindFound = false

			found, err = hasDisk.Run("fake-disk-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(diskService.FindCalled).To(BeTrue())
		})

		It("returns an error if diskService find call returns an error", func() {
			diskService.FindErr = errors.New("fake-vm-service-error")

			_, err = hasDisk.Run("fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(diskService.FindCalled).To(BeTrue())
		})
	})
})
