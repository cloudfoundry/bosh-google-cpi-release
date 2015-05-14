package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"

	imagefakes "github.com/frodenas/bosh-google-cpi/google/image_service/fakes"
)

var _ = Describe("CreateStemcell", func() {
	var (
		err         error
		stemcellCID StemcellCID
		cloudProps  StemcellCloudProperties

		stemcellService *imagefakes.FakeImageService
		createStemcell  CreateStemcell
	)

	BeforeEach(func() {
		stemcellService = &imagefakes.FakeImageService{}
		createStemcell = NewCreateStemcell(stemcellService)
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
				Expect(stemcellService.CreateFromTarballCalled).To(BeFalse())
				Expect(stemcellService.CreateFromURLCalled).To(BeFalse())
			})
		})

		Context("from a source url", func() {
			BeforeEach(func() {
				cloudProps.SourceURL = "fake-source-url"
				stemcellService.CreateFromURLID = "fake-stemcell-id"
			})

			It("creates the stemcell", func() {
				stemcellCID, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(stemcellService.CreateFromURLCalled).To(BeTrue())
				Expect(stemcellService.CreateFromTarballCalled).To(BeFalse())
				Expect(stemcellCID).To(Equal(StemcellCID("fake-stemcell-id")))
				Expect(stemcellService.CreateFromURLSourceURL).To(Equal("fake-source-url"))
				Expect(stemcellService.CreateFromURLDescription).To(Equal("fake-stemcell-name/fake-stemcell-version"))
			})

			It("returns an error if stemcellService create from tarball call returns an error", func() {
				stemcellService.CreateFromURLErr = errors.New("fake-stemcell-service-error")

				_, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-stemcell-service-error"))
				Expect(stemcellService.CreateFromURLCalled).To(BeTrue())
				Expect(stemcellService.CreateFromTarballCalled).To(BeFalse())
			})
		})

		Context("from a stemcell tarball", func() {
			BeforeEach(func() {
				stemcellService.CreateFromTarballID = "fake-stemcell-id"
			})

			It("creates the stemcell", func() {
				stemcellCID, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(stemcellService.CreateFromTarballCalled).To(BeTrue())
				Expect(stemcellService.CreateFromURLCalled).To(BeFalse())
				Expect(stemcellCID).To(Equal(StemcellCID("fake-stemcell-id")))
				Expect(stemcellService.CreateFromTarballImagePath).To(Equal("fake-stemcell-tarball"))
				Expect(stemcellService.CreateFromTarballDescription).To(Equal("fake-stemcell-name/fake-stemcell-version"))
			})

			It("returns an error if stemcellService create from tarball call returns an error", func() {
				stemcellService.CreateFromTarballErr = errors.New("fake-stemcell-service-error")

				_, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-stemcell-service-error"))
				Expect(stemcellService.CreateFromTarballCalled).To(BeTrue())
				Expect(stemcellService.CreateFromURLCalled).To(BeFalse())
			})
		})
	})
})
