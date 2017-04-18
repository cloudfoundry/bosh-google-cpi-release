package config


import (
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
        var (
                config Config
        )

	Describe("UserAgent", func() {
                It("returns correct user agent string with release, without prefix", func() {
                        config.UserAgentPrefix = ""
                        cpiRelease = "0.0.1"

                        userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("bosh-google-cpi/0.0.1"))
                })
                It("returns correct user agent string with release, with prefix", func() {
                        config.UserAgentPrefix = "Kubo/0.0.2"
                        cpiRelease = "0.0.1"

                        userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("Kubo/0.0.2 bosh-google-cpi/0.0.1"))
                })
                It("returns correct user agent string without release, with prefix", func() {
                        config.UserAgentPrefix = "Kubo/0.0.2"
                        cpiRelease = ""

                        userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("Kubo/0.0.2 bosh-google-cpi/dev"))
                })
                It("returns correct user agent string without release, without prefix", func() {
                        config.UserAgentPrefix = ""
                        cpiRelease = ""

                        userAgent := config.GetUserAgent()
                        Expect(userAgent).To(Equal("bosh-google-cpi/dev"))
                })
        })
})
