package registry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/client"
)

var _ = Describe("ClientOptions", func() {
	var (
		options ClientOptions

		validOptions = ClientOptions{
			Protocol: "http",
			Host:     "fake-host",
			Port:     25777,
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

		It("returns error if Protocol is empty", func() {
			options.Protocol = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Protocol"))
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

var _ = Describe("ClientTLSOptions", func() {
	var (
		options ClientOptions

		validOptions = ClientOptions{
			Protocol: "https",
			Host:     "fake-host",
			Port:     25777,
			Username: "fake-username",
			Password: "fake-password",
			TLS: ClientTLSOptions{
				CertFile: "fake-certificate",
				KeyFile:  "fake-key",
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

		It("returns error if CertFile is empty", func() {
			options.TLS.CertFile = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty CertFile"))
		})

		It("returns error if KeyFile is empty", func() {
			options.TLS.KeyFile = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty KeyFile"))
		})
	})
})
