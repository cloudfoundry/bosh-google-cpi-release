package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"
	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk"
	diskfakes "bosh-google-cpi/google/disk/fakes"
	"bosh-google-cpi/google/disktype"
	disktypefakes "bosh-google-cpi/google/disktype/fakes"
	"bosh-google-cpi/google/snapshot"
	snapshotfakes "bosh-google-cpi/google/snapshot/fakes"
)

var _ = Describe("UpdateDisk", func() {
	var (
		err    error
		result interface{}

		diskService     *diskfakes.FakeDiskService
		diskTypeService *disktypefakes.FakeDiskTypeService
		snapshotService *snapshotfakes.FakeSnapshotService

		updateDisk UpdateDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		diskTypeService = &disktypefakes.FakeDiskTypeService{}
		snapshotService = &snapshotfakes.FakeSnapshotService{}
		updateDisk = NewUpdateDisk(diskService, diskTypeService, snapshotService)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			diskService.FindFound = true
			diskService.FindDisk = disk.Disk{
				Name:   "fake-disk-id",
				Zone:   "fake-zone",
				SizeGb: 10,
				Type:   "https://googleapis.com/zones/fake-zone/diskTypes/pd-ssd",
			}
		})

		Context("when disk is not found", func() {
			BeforeEach(func() {
				diskService.FindFound = false
			})

			It("returns a DiskNotFoundError", func() {
				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(api.NewDiskNotFoundError("fake-disk-id", false).Error()))
			})
		})

		Context("when finding disk returns an error", func() {
			BeforeEach(func() {
				diskService.FindErr = errors.New("fake-find-error")
			})

			It("returns an error", func() {
				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-find-error"))
			})
		})

		Context("when type and size are unchanged (no-op)", func() {
			It("returns the existing disk CID without any modifications", func() {
				result, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{})
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("fake-disk-id"))
				Expect(diskService.ResizeCalled).To(BeFalse())
				Expect(snapshotService.CreateCalled).To(BeFalse())
			})

			It("treats same type as no-op even when self-link URL formats differ", func() {
				// Existing disk has a full self-link from the GCP Disks API
				diskService.FindDisk = disk.Disk{
					Name:   "fake-disk-id",
					Zone:   "fake-zone",
					SizeGb: 10,
					Type:   "https://www.googleapis.com/compute/v1/projects/my-proj/zones/us-east1-b/diskTypes/pd-ssd",
				}
				// DiskType.Find returns a self-link from the DiskTypes API (same format in practice,
				// but we test with a slightly different prefix to prove ResourceSplitter handles it)
				diskTypeService.FindFound = true
				diskTypeService.FindDiskType = disktype.DiskType{
					SelfLink: "https://compute.googleapis.com/compute/v1/projects/my-proj/zones/us-east1-b/diskTypes/pd-ssd",
				}

				result, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "pd-ssd"})
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("fake-disk-id"))
				// Should NOT trigger snapshot cycle — same type name "pd-ssd"
				Expect(snapshotService.CreateCalled).To(BeFalse())
				Expect(diskService.ResizeCalled).To(BeFalse())
			})

			It("treats explicitly specified current type as no-op", func() {
				diskTypeService.FindFound = true
				diskTypeService.FindDiskType = disktype.DiskType{
					SelfLink: "https://googleapis.com/zones/fake-zone/diskTypes/pd-ssd",
				}

				result, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "pd-ssd"})
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("fake-disk-id"))
				Expect(snapshotService.CreateCalled).To(BeFalse())
			})
		})

		Context("when only size increases (same type)", func() {
			It("resizes in-place and returns the existing disk CID", func() {
				result, err = updateDisk.Run("fake-disk-id", 20480, DiskCloudProperties{})
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("fake-disk-id"))
				Expect(diskService.ResizeCalled).To(BeTrue())
				Expect(snapshotService.CreateCalled).To(BeFalse())
			})

			It("returns an error if resize fails", func() {
				diskService.ResizeErr = errors.New("fake-resize-error")

				_, err = updateDisk.Run("fake-disk-id", 20480, DiskCloudProperties{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-resize-error"))
			})
		})

		Context("when requested size is smaller than current", func() {
			It("returns an error", func() {
				_, err = updateDisk.Run("fake-disk-id", 5120, DiskCloudProperties{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("smaller than current size"))
			})
		})

		Context("when disk type changes", func() {
			BeforeEach(func() {
				diskTypeService.FindFound = true
				diskTypeService.FindDiskType = disktype.DiskType{
					SelfLink: "https://googleapis.com/zones/fake-zone/diskTypes/hyperdisk-balanced",
				}

				snapshotService.CreateID = "fake-snapshot-id"
				snapshotService.FindFound = true
				snapshotService.FindSnapshot = snapshot.Snapshot{
					SelfLink: "https://googleapis.com/snapshots/fake-snapshot-id",
				}

				diskService.CreateFromSnapshotID = "new-disk-id"
			})

			It("snapshots, creates new disk from snapshot, deletes old disk, and returns new CID", func() {
				result, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "hyperdisk-balanced"})
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("new-disk-id"))
				Expect(snapshotService.CreateCalled).To(BeTrue())
				Expect(snapshotService.FindCalled).To(BeTrue())
				Expect(diskService.CreateFromSnapshotCalled).To(BeTrue())
				Expect(diskService.CreateFromSnapshotSnapshotSelfLink).To(Equal("https://googleapis.com/snapshots/fake-snapshot-id"))
				Expect(diskService.CreateFromSnapshotDiskType).To(Equal("https://googleapis.com/zones/fake-zone/diskTypes/hyperdisk-balanced"))
				Expect(diskService.DeleteCalled).To(BeTrue())
				Expect(snapshotService.DeleteCalled).To(BeTrue())
			})

			It("returns an error if diskTypeService find fails", func() {
				diskTypeService.FindErr = errors.New("fake-disktype-error")

				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "hyperdisk-balanced"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-disktype-error"))
			})

			It("returns an error if disk type is not found", func() {
				diskTypeService.FindFound = false

				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "nonexistent"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})

			It("returns an error if snapshot creation fails", func() {
				snapshotService.CreateErr = errors.New("fake-snapshot-create-error")

				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "hyperdisk-balanced"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-snapshot-create-error"))
			})

			It("returns an error and cleans up snapshot if snapshot find fails", func() {
				snapshotService.FindErr = errors.New("fake-snapshot-find-error")

				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "hyperdisk-balanced"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-snapshot-find-error"))
				Expect(snapshotService.DeleteCalled).To(BeTrue())
			})

			It("returns an error and cleans up snapshot if snapshot not found after creation", func() {
				snapshotService.FindFound = false

				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "hyperdisk-balanced"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found after creation"))
				Expect(snapshotService.DeleteCalled).To(BeTrue())
			})

			It("returns an error and cleans up snapshot if CreateFromSnapshot fails", func() {
				diskService.CreateFromSnapshotErr = errors.New("fake-create-from-snapshot-error")

				_, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "hyperdisk-balanced"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-create-from-snapshot-error"))
				Expect(snapshotService.DeleteCalled).To(BeTrue())
			})

			It("returns new disk CID with error if old disk deletion fails", func() {
				diskService.DeleteErr = errors.New("fake-delete-error")

				result, err = updateDisk.Run("fake-disk-id", 10240, DiskCloudProperties{DiskType: "hyperdisk-balanced"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-delete-error"))
				Expect(result).To(Equal("new-disk-id"))
			})
		})
	})
})
