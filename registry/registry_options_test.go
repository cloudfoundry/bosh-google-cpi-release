package registry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/registry"
)

var _ = Describe("RegistryOptions", func() {
	var (
		options RegistryOptions

		validOptions = RegistryOptions{
			Schema:   "http",
			Host:     "fake-host",
			Port:     5555,
			Username: "fake-username",
			Password: "fake-password",
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

		It("returns error if Schema is empty", func() {
			options.Schema = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Schema"))
		})

		It("returns error if Host is empty", func() {
			options.Host = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Host"))
		})

		It("returns error if Port is empty", func() {
			options.Port = 0

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Port"))
		})

		It("returns error if Username is empty", func() {
			options.Username = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Username"))
		})

		It("returns error if Password is empty", func() {
			options.Password = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Password"))
		})
	})
})
