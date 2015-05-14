package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"

	diskfakes "github.com/frodenas/bosh-google-cpi/google/disk_service/fakes"
	instancefakes "github.com/frodenas/bosh-google-cpi/google/instance_service/fakes"
	registryfakes "github.com/frodenas/bosh-registry/client/fakes"

	"github.com/frodenas/bosh-google-cpi/api"
	"google.golang.org/api/compute/v1"
)

var _ = Describe("AttachDisk", func() {
	var (
		err            error
		diskService    *diskfakes.FakeDiskService
		vmService      *instancefakes.FakeInstanceService
		registryClient *registryfakes.FakeClient
		attachDisk     AttachDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		vmService = &instancefakes.FakeInstanceService{}
		registryClient = &registryfakes.FakeClient{}
		attachDisk = NewAttachDisk(diskService, vmService, registryClient)
	})

	Describe("Run", func() {
		It("attaches the disk", func() {
			diskService.FindFound = true
			diskService.FindDisk = &compute.Disk{SelfLink: "fake-self-link"}

			_, err = attachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(vmService.AttachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})

		It("returns an error if disk is not found", func() {
			diskService.FindFound = false

			_, err = attachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(api.NewDiskNotFoundError("fake-disk-id", false).Error()))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(vmService.AttachDiskCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if diskService find call returns an error", func() {
			diskService.FindErr = errors.New("fake-disk-service-error")

			_, err = attachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-disk-service-error"))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(vmService.AttachDiskCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if vmService attach disk call returns an error", func() {
			diskService.FindFound = true
			diskService.FindDisk = &compute.Disk{SelfLink: "fake-self-link"}
			vmService.AttachDiskErr = errors.New("fake-vm-service-error")

			_, err = attachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(vmService.AttachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient fetch call returns an error", func() {
			diskService.FindFound = true
			diskService.FindDisk = &compute.Disk{SelfLink: "fake-self-link"}
			registryClient.FetchErr = errors.New("fake-registry-client-error")

			_, err = attachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(vmService.AttachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			diskService.FindFound = true
			diskService.FindDisk = &compute.Disk{SelfLink: "fake-self-link"}
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = attachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(vmService.AttachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})
	})
})
