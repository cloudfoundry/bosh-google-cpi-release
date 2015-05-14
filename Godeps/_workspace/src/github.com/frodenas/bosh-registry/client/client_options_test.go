package registry_test

import (
	"fmt"

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

	BeforeEach(func() {
		options = validOptions
	})

	Describe("Endpoint", func() {
		It("returns the BOSH Registry endpoint", func() {
			endpoint := options.Endpoint()
			Expect(endpoint).To(Equal(fmt.Sprintf("%s://%s:%d", validOptions.Protocol, validOptions.Host, validOptions.Port)))
		})
	})

	Describe("EndpointWithCredentials", func() {
		It("returns the BOSH Registry endpoint with credentials", func() {
			endpoint := options.EndpointWithCredentials()
			Expect(endpoint).To(Equal(fmt.Sprintf("%s://%s:%s@%s:%d", validOptions.Protocol, validOptions.Username, validOptions.Password, validOptions.Host, validOptions.Port)))
		})
	})

	Describe("Validate", func() {
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

	BeforeEach(func() {
		options = validOptions
	})

	Describe("Validate", func() {
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
