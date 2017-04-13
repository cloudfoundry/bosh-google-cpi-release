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
		It("returns correct user agent string with release, without prefix", func() {
			config.UserAgentPrefix = ""
			CpiRelease = "0.0.1"
			
			userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("bosh-google-cpi/0.0.1"))
		})
		It("returns correct user agent string with release, with prefix", func() {
			config.UserAgentPrefix = "Kubo/0.0.2"
			CpiRelease = "0.0.1"

			userAgent := config.GetUserAgent()
			Expect(userAgent).To(Equal("Kubo/0.0.2 bosh-google-cpi/0.0.1"))
		})
		It("returns correct user agent string without release, with prefix", func() {
                        config.UserAgentPrefix = "Kubo/0.0.2"
                        CpiRelease = ""

                        userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("Kubo/0.0.2 bosh-google-cpi/dev"))
                })
		It("returns correct user agent string without release, without prefix", func() {
                        config.UserAgentPrefix = ""
                        CpiRelease = ""

                        userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("bosh-google-cpi/dev"))
                })
	})
})
