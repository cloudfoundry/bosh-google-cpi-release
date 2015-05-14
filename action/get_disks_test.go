package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"

	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	instancefakes "github.com/frodenas/bosh-google-cpi/google/instance_service/fakes"
)

var _ = Describe("GetDisks", func() {
	var (
		err       error
		vmService *instancefakes.FakeInstanceService
		getDisks  GetDisks
		disks     ginstance.InstanceAttachedDisks
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		getDisks = NewGetDisks(vmService)
	})

	Describe("Run", func() {
		It("returns the list of attached disks", func() {
			attachedDisksList := []string{"fake-disk-1", "fake-disk-2"}
			vmService.AttachedDisksList = ginstance.InstanceAttachedDisks(attachedDisksList)

			disks, err = getDisks.Run("fake-vm-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.AttachedDisksCalled).To(BeTrue())
			Expect(disks).To(Equal(ginstance.InstanceAttachedDisks(attachedDisksList)))
		})

		Context("when there are not any attached disk", func() {
			It("returns an empty array", func() {
				disks, err = getDisks.Run("fake-vm-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.AttachedDisksCalled).To(BeTrue())
				Expect(disks).To(BeEmpty())
			})
		})

		It("returns an error if vmService attached disks call returns an error", func() {
			vmService.AttachedDisksErr = errors.New("fake-vm-service-error")

			_, err = getDisks.Run("fake-vm-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.AttachedDisksCalled).To(BeTrue())
		})
	})
})
