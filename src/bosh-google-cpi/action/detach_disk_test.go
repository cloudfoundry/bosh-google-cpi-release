package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	instancefakes "bosh-google-cpi/google/instance_service/fakes"
	registryfakes "github.com/frodenas/bosh-registry/client/fakes"

	"github.com/frodenas/bosh-registry/client"
)

var _ = Describe("DetachDisk", func() {
	var (
		err                   error
		expectedAgentSettings registry.AgentSettings

		vmService      *instancefakes.FakeInstanceService
		registryClient *registryfakes.FakeClient

		detachDisk DetachDisk
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		registryClient = &registryfakes.FakeClient{}
		detachDisk = NewDetachDisk(vmService, registryClient)
		registryClient.FetchSettings = registry.AgentSettings{
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
		expectedAgentSettings = registry.AgentSettings{
			Disks: registry.DisksSettings{
				Persistent: map[string]registry.PersistentSettings{},
			},
		}
	})

	Describe("Run", func() {
		It("detaches the disk", func() {
			_, err = detachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.DetachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
		})

		It("returns an error if vmService detach disk call returns an error", func() {
			vmService.DetachDiskErr = errors.New("fake-vm-service-error")

			_, err = detachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.DetachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient fetch call returns an error", func() {
			registryClient.FetchErr = errors.New("fake-registry-client-error")

			_, err = detachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(vmService.DetachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = detachDisk.Run("fake-vm-id", "fake-disk-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(vmService.DetachDiskCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})
	})
})
