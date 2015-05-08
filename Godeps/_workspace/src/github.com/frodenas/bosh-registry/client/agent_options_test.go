package registry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/client"
)

var _ = Describe("AgentOptions", func() {
	var (
		options AgentOptions

		validOptions = AgentOptions{
			Mbus: "fake-mbus",
			Ntp:  []string{},
			Blobstore: BlobstoreOptions{
				Type: "fake-blobstore-type",
			},
		}
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			options = validOptions
		})

		It("does not return error if all fields are valid", func() {
			err := options.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if Mbus is empty", func() {
			options.Mbus = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Mbus"))
		})

		It("returns error if blobstore options are not valid", func() {
			options.Blobstore = BlobstoreOptions{}

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Blobstore configuration"))
		})
	})
})

var _ = Describe("BlobstoreOptions", func() {
	var (
		options BlobstoreOptions

		validOptions = BlobstoreOptions{
			Type:    "fake-type",
			Options: map[string]interface{}{"fake-key": "fake-value"},
		}
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			options = validOptions
		})

		It("does not return error if all fields are valid", func() {
			err := options.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if Type is empty", func() {
			options.Type = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide non-empty Type"))
		})
	})
})
