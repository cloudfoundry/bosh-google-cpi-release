package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	diskfakes "bosh-google-cpi/google/disk_service/fakes"
	instancefakes "bosh-google-cpi/google/instance_service/fakes"

	registryfakes "bosh-google-cpi/registry/fakes"

	"bosh-google-cpi/api"
	"bosh-google-cpi/google/disk_service"

	"bosh-google-cpi/registry"
)

var _ = Describe("AttachDisk", func() {
	var (
		err                   error
		expectedAgentSettings registry.AgentSettings

		diskService    *diskfakes.FakeDiskService
		vmService      *instancefakes.FakeInstanceService
		registryClient *registryfakes.FakeClient

		attachDisk AttachDisk
	)

	BeforeEach(func() {
		diskService = &diskfakes.FakeDiskService{}
		vmService = &instancefakes.FakeInstanceService{}
		registryClient = &registryfakes.FakeClient{}
		attachDisk = NewAttachDisk(diskService, vmService, registryClient)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			diskService.FindFound = true
			diskService.FindDisk = disk.Disk{SelfLink: "fake-self-link"}
			vmService.AttachDiskDeviceName = "fake-disk-device-name"
			vmService.AttachDiskDevicePath = "fake-disk-device-path"
			registryClient.FetchSettings = registry.AgentSettings{}
			expectedAgentSettings = registry.AgentSettings{
				Disks: registry.DisksSettings{
					Persistent: map[string]registry.PersistentSettings{
						"fake-disk-id": {
							ID:       "fake-disk-id",
							VolumeID: "fake-disk-device-name",
							Path:     "fake-disk-device-path",
						},
					},
				},
			}
		})

		It("attaches the disk", func() {
			_, err = attachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(diskService.FindCalled).To(BeTrue())
			Expect(vmService.AttachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
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

		It("returns an error if vmService attach disk call returns an error", func() {
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
