package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"

	imagefakes "github.com/frodenas/bosh-google-cpi/google/image_service/fakes"
)

var _ = Describe("DeleteStemcell", func() {
	var (
		err error

		stemcellService *imagefakes.FakeImageService

		deleteStemcell DeleteStemcell
	)

	BeforeEach(func() {
		stemcellService = &imagefakes.FakeImageService{}
		deleteStemcell = NewDeleteStemcell(stemcellService)
	})

	Describe("Run", func() {
		It("deletes the stemcell", func() {
			_, err = deleteStemcell.Run("fake-stemcell-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellService.DeleteCalled).To(BeTrue())
		})

		It("returns an error if stemcellService delete call returns an error", func() {
			stemcellService.DeleteErr = errors.New("fake-stemcell-service-error")

			_, err = deleteStemcell.Run("fake-stemcell-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-stemcell-service-error"))
			Expect(stemcellService.DeleteCalled).To(BeTrue())
		})
	})
})
