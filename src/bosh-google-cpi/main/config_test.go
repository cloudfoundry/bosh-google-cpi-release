package main_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "bosh-google-cpi/main"

	bgcaction "bosh-google-cpi/action"
	bgcconfig "bosh-google-cpi/google/config"

	"github.com/frodenas/bosh-registry/client"
)

var validGoogleConfig = bgcconfig.Config{
	Project:     "fake-project",
	DefaultZone: "fake-default-zone",
}

var validActionsOptions = bgcaction.ConcreteFactoryOptions{
	Agent: registry.AgentOptions{
		Mbus: "fake-mbus",
		Ntp:  []string{},
		Blobstore: registry.BlobstoreOptions{
			Type: "fake-blobstore-type",
		},
	},
	Registry: registry.ClientOptions{
		Protocol: "http",
		Host:     "fake-host",
		Port:     5555,
		Username: "fake-username",
		Password: "fake-password",
	},
}

var validConfig = Config{
	Google:  validGoogleConfig,
	Actions: validActionsOptions,
}

var _ = Describe("NewConfigFromPath", func() {
	var (
		fs *fakesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
	})

	It("returns error if config is empty", func() {
		_, err := NewConfigFromPath("", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Must provide a config file"))
	})

	It("returns error if config is not valid", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Validating config"))
	})

	It("returns error if file contains invalid json", func() {
		err := fs.WriteFileString("/config.json", "-")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Unmarshalling config"))
	})

	It("returns error if file cannot be read", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		fs.ReadFileError = errors.New("fake-read-err")

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-read-err"))
	})
})

var _ = Describe("Config", func() {
	var (
		config Config
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validConfig
		})

		It("does not return error if all google and actions sections are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if goole section is not valid", func() {
			config.Google = bgcconfig.Config{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Google configuration"))
		})

		It("returns error if actions section is not valid", func() {
			config.Actions.Agent = registry.AgentOptions{}
			config.Actions.Registry = registry.ClientOptions{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Actions configuration"))
		})
	})
})
