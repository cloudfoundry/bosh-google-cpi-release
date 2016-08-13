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
})
