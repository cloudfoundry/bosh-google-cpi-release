package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	diskfakes "bosh-google-cpi/google/disk_service/fakes"
)

var _ = Describe("ResizeDisk", func() {
	var (
		err error

		diskService *diskfakes.FakeDiskService

		resizeDisk ResizeDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		resizeDisk = NewResizeDisk(diskService)
	})
	//TODO: write tests for:
	// cannot resize to smaller disk (this should be handled already by the director call)
	// skip resize when values are the same
	// cannot resize as disk is still attached.
	Describe("Run", func() {
		It("resize the disk", func() {
			_, err = resizeDisk.Run("fake-disk-id", 123)
			Expect(err).NotTo(HaveOccurred())
			Expect(diskService.ResizeCalled).To(BeTrue())
		})

		It("returns an error if diskService resize call returns an error", func() {
			diskService.ResizeErr = errors.New("fake-disk-service-error")

			_, err = resizeDisk.Run("fake-disk-id", 123)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-disk-service-error"))
			Expect(diskService.ResizeCalled).To(BeTrue())
		})
	})
})
