package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	diskfakes "bosh-google-cpi/google/disk_service/fakes"
	disktypefakes "bosh-google-cpi/google/disk_type_service/fakes"
	instancefakes "bosh-google-cpi/google/instance_service/fakes"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_type_service"

	"google.golang.org/api/compute/v1"
)

var _ = Describe("CreateDisk", func() {
	var (
		err        error
		diskCID    DiskCID
		vmCID      VMCID
		cloudProps DiskCloudProperties

		diskService     *diskfakes.FakeDiskService
		diskTypeService *disktypefakes.FakeDiskTypeService
		vmService       *instancefakes.FakeInstanceService

		createDisk CreateDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		diskTypeService = &disktypefakes.FakeDiskTypeService{}
		vmService = &instancefakes.FakeInstanceService{}
		createDisk = NewCreateDisk(diskService, diskTypeService, vmService, "fake-default-zone")
	})

	Describe("Run", func() {
		BeforeEach(func() {
			vmCID = ""
			cloudProps = DiskCloudProperties{}
			diskService.CreateID = "fake-disk-id"
		})

		It("creates the disk", func() {
			diskCID, err = createDisk.Run(32768, cloudProps, vmCID)
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(diskService.CreateCalled).To(BeTrue())
			Expect(diskService.CreateSize).To(Equal(32))
			Expect(diskService.CreateDiskType).To(BeEmpty())
			Expect(diskService.CreateZone).To(Equal("fake-default-zone"))
			Expect(diskCID).To(Equal(DiskCID("fake-disk-id")))
		})

		It("returns an error if diskService create call returns an error", func() {
			diskService.CreateErr = errors.New("fake-disk-service-error")

			_, err = createDisk.Run(32768, cloudProps, vmCID)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-disk-service-error"))
			Expect(vmService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(diskService.CreateCalled).To(BeTrue())
		})

		Context("when vmCID is set", func() {
			BeforeEach(func() {
				vmCID = VMCID("fake-vm-cid")
				vmService.FindFound = true
				vmService.FindInstance = &compute.Instance{Zone: "fake-instance-zone"}
			})

			It("creates the disk at the vm zone", func() {
				diskCID, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(diskService.CreateCalled).To(BeTrue())
				Expect(diskService.CreateSize).To(Equal(32))
				Expect(diskService.CreateDiskType).To(BeEmpty())
				Expect(diskService.CreateZone).To(Equal("fake-instance-zone"))
				Expect(diskCID).To(Equal(DiskCID("fake-disk-id")))
			})

			It("returns an error if vmService find call returns an error", func() {
				vmService.FindErr = errors.New("fake-instance-service-error")

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-instance-service-error"))
				Expect(vmService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(diskService.CreateCalled).To(BeFalse())
			})

			It("returns an error if instance is not found", func() {
				vmService.FindFound = false

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(api.NewVMNotFoundError(string(vmCID)).Error()))
				Expect(vmService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(diskService.CreateCalled).To(BeFalse())
			})
		})

		Context("when disk type is set", func() {
			BeforeEach(func() {
				cloudProps = DiskCloudProperties{DiskType: "fake-disk-type"}
				diskTypeService.FindFound = true
				diskTypeService.FindDiskType = disktype.DiskType{SelfLink: "fake-disk-type-self-link"}
			})

			It("creates the disk using the appropiate disk type", func() {
				diskCID, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(diskService.CreateCalled).To(BeTrue())
				Expect(diskService.CreateSize).To(Equal(32))
				Expect(diskService.CreateDiskType).To(Equal("fake-disk-type-self-link"))
				Expect(diskService.CreateZone).To(Equal("fake-default-zone"))
				Expect(diskCID).To(Equal(DiskCID("fake-disk-id")))
			})

			It("returns an error if diskTypeService find call returns an error", func() {
				diskTypeService.FindErr = errors.New("fake-disk-type-service-error")

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-disk-type-service-error"))
				Expect(vmService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(diskService.CreateCalled).To(BeFalse())
			})

			It("returns an error if disk type is not found", func() {
				diskTypeService.FindFound = false

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Disk Type 'fake-disk-type' does not exists"))
				Expect(vmService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(diskService.CreateCalled).To(BeFalse())
			})
		})
	})

})
