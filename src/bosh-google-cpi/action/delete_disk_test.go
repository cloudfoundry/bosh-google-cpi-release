package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	diskfakes "bosh-google-cpi/google/disk_service/fakes"
)

var _ = Describe("DeleteDisk", func() {
	var (
		err error

		diskService *diskfakes.FakeDiskService

		deleteDisk DeleteDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		deleteDisk = NewDeleteDisk(diskService)
	})

	Describe("Run", func() {
		It("deletes the disk", func() {
			_, err = deleteDisk.Run("fake-disk-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(diskService.DeleteCalled).To(BeTrue())
		})

		It("returns an error if diskService delete call returns an error", func() {
			diskService.DeleteErr = errors.New("fake-disk-service-error")

			_, err = deleteDisk.Run("fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-disk-service-error"))
			Expect(diskService.DeleteCalled).To(BeTrue())
		})
	})
})
