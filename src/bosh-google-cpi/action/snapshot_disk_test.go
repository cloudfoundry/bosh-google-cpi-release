package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	diskfakes "bosh-google-cpi/google/disk_service/fakes"
	snapshotfakes "bosh-google-cpi/google/snapshot_service/fakes"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_service"
)

var _ = Describe("SnapshotDisk", func() {
	var (
		err        error
		metadata   SnapshotMetadata
		snapshotID SnapshotCID

		diskService     *diskfakes.FakeDiskService
		snapshotService *snapshotfakes.FakeSnapshotService

		snapshotDisk SnapshotDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		snapshotService = &snapshotfakes.FakeSnapshotService{}
		snapshotDisk = NewSnapshotDisk(snapshotService, diskService)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			diskService.FindFound = true
			diskService.FindDisk = disk.Disk{Zone: "fake-disk-zone"}
			snapshotService.CreateID = "fake-snapshot-id"
			metadata = SnapshotMetadata{Deployment: "fake-deployment", Job: "fake-job", Index: "fake-index"}
		})

		Context("creates a snaphot", func() {
			It("with the proper description", func() {
				snapshotID, err = snapshotDisk.Run("fake-disk-id", metadata)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeTrue())
				Expect(snapshotService.CreateCalled).To(BeTrue())
				Expect(snapshotService.CreateDiskID).To(Equal("fake-disk-id"))
				Expect(snapshotService.CreateDescription).To(Equal("fake-deployment/fake-job/fake-index"))
				Expect(snapshotService.CreateZone).To(Equal("fake-disk-zone"))
				Expect(snapshotID).To(Equal(SnapshotCID("fake-snapshot-id")))
			})

			Context("when metadata is empty", func() {
				BeforeEach(func() {
					metadata = SnapshotMetadata{}
				})

				It("with an empty description", func() {
					snapshotID, err = snapshotDisk.Run("fake-disk-id", metadata)
					Expect(err).NotTo(HaveOccurred())
					Expect(diskService.FindCalled).To(BeTrue())
					Expect(snapshotService.CreateCalled).To(BeTrue())
					Expect(snapshotService.CreateDiskID).To(Equal("fake-disk-id"))
					Expect(snapshotService.CreateDescription).To(BeEmpty())
					Expect(snapshotService.CreateZone).To(Equal("fake-disk-zone"))
					Expect(snapshotID).To(Equal(SnapshotCID("fake-snapshot-id")))
				})
			})
		})

		It("returns an error if diskService find call returns an error", func() {
			diskService.FindErr = errors.New("fake-disk-service-error")

			_, err = snapshotDisk.Run("fake-disk-id", metadata)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-disk-service-error"))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(snapshotService.CreateCalled).To(BeFalse())
		})

		It("returns an error if disk is not found", func() {
			diskService.FindFound = false

			_, err = snapshotDisk.Run("fake-disk-id", metadata)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(api.NewDiskNotFoundError("fake-disk-id", false).Error()))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(snapshotService.CreateCalled).To(BeFalse())
		})

		It("returns an error if snapshotService create call returns an error", func() {
			snapshotService.CreateErr = errors.New("fake-snapshot-service-error")

			_, err = snapshotDisk.Run("fake-disk-id", metadata)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-snapshot-service-error"))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(snapshotService.CreateCalled).To(BeTrue())
		})
	})
})
