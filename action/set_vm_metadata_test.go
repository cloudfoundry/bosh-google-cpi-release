package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"

	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	instancefakes "github.com/frodenas/bosh-google-cpi/google/instance_service/fakes"
)

var _ = Describe("SetVMMetadata", func() {
	var (
		err           error
		vmService     *instancefakes.FakeInstanceService
		setVMMetadata SetVMMetadata
		vmMetadata    VMMetadata
	)

	BeforeEach(func() {
		vmMetadata = map[string]interface{}{
			"deployment": "fake-deployment",
			"job":        "fake-job",
			"index":      "fake-index",
		}
		vmService = &instancefakes.FakeInstanceService{}
		setVMMetadata = NewSetVMMetadata(vmService)
	})

	Describe("Run", func() {
		It("set the vm metadata", func() {
			_, err = setVMMetadata.Run("fake-vm-id", vmMetadata)
			Expect(err).NotTo(HaveOccurred())
			Expect(vmService.SetMetadataCalled).To(BeTrue())
			Expect(vmService.SetMetadataVMMetadata).To(Equal(ginstance.InstanceMetadata(vmMetadata)))
		})

		It("returns an error if vmService set metadata call returns an error", func() {
			vmService.SetMetadataErr = errors.New("fake-vm-service-error")

			_, err = setVMMetadata.Run("fake-vm-id", vmMetadata)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(vmService.SetMetadataCalled).To(BeTrue())
		})
	})
})
