package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/google/config"
)

var validConfig = Config{
	Project: "fake-project",
}

var _ = Describe("Config", func() {
	var (
		config Config
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validConfig
		})

		It("does not return error if all fields are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if Project is empty", func() {
			config.Project = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Project"))
		})
	})
	Describe("UserAgent", func() {
		It("returns a valid user agent string without external user agent", func() {
			config.UserAgentPrefix = ""
			
			userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("bosh-google-cpi/dev"))
		})
		It("returns a valid user agent string with external user agent", func() {
			config.UserAgentPrefix = "Kubo/0.0.2"

			userAgent := config.GetUserAgent()
			Expect(userAgent).To(Equal("Kubo/0.0.2 bosh-google-cpi/dev"))
		})
	})
})
