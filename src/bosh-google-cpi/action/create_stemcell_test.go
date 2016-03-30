package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	imagefakes "bosh-google-cpi/google/image_service/fakes"
)

var _ = Describe("CreateStemcell", func() {
	var (
		err         error
		stemcellCID StemcellCID
		cloudProps  StemcellCloudProperties

		imageService   *imagefakes.FakeImageService
		createStemcell CreateStemcell
	)

	BeforeEach(func() {
		imageService = &imagefakes.FakeImageService{}
		createStemcell = NewCreateStemcell(imageService)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			cloudProps = StemcellCloudProperties{
				Name:           "fake-stemcell-name",
				Version:        "fake-stemcell-version",
				Infrastructure: "google",
			}
		})

		Context("when infrastructure is not google", func() {
			BeforeEach(func() {
				cloudProps.Infrastructure = "fake-insfrastructure"
			})

			It("returns an error", func() {
				_, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid 'fake-insfrastructure' infrastructure"))
				Expect(imageService.CreateFromTarballCalled).To(BeFalse())
				Expect(imageService.CreateFromURLCalled).To(BeFalse())
			})
		})

		Context("from a source url", func() {
			BeforeEach(func() {
				cloudProps.SourceURL = "fake-source-url"
				imageService.CreateFromURLID = "fake-stemcell-id"
			})

			It("creates the stemcell", func() {
				stemcellCID, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(imageService.CreateFromURLCalled).To(BeTrue())
				Expect(imageService.CreateFromTarballCalled).To(BeFalse())
				Expect(stemcellCID).To(Equal(StemcellCID("fake-stemcell-id")))
				Expect(imageService.CreateFromURLSourceURL).To(Equal("fake-source-url"))
				Expect(imageService.CreateFromURLDescription).To(Equal("fake-stemcell-name/fake-stemcell-version"))
			})

			It("returns an error if imageService create from tarball call returns an error", func() {
				imageService.CreateFromURLErr = errors.New("fake-image-service-error")

				_, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-image-service-error"))
				Expect(imageService.CreateFromURLCalled).To(BeTrue())
				Expect(imageService.CreateFromTarballCalled).To(BeFalse())
			})
		})

		Context("from a stemcell tarball", func() {
			BeforeEach(func() {
				imageService.CreateFromTarballID = "fake-stemcell-id"
			})

			It("creates the stemcell", func() {
				stemcellCID, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(imageService.CreateFromTarballCalled).To(BeTrue())
				Expect(imageService.CreateFromURLCalled).To(BeFalse())
				Expect(stemcellCID).To(Equal(StemcellCID("fake-stemcell-id")))
				Expect(imageService.CreateFromTarballImagePath).To(Equal("fake-stemcell-tarball"))
				Expect(imageService.CreateFromTarballDescription).To(Equal("fake-stemcell-name/fake-stemcell-version"))
			})

			It("returns an error if imageService create from tarball call returns an error", func() {
				imageService.CreateFromTarballErr = errors.New("fake-stemcell-service-error")

				_, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-stemcell-service-error"))
				Expect(imageService.CreateFromTarballCalled).To(BeTrue())
				Expect(imageService.CreateFromURLCalled).To(BeFalse())
			})
		})
	})
})
