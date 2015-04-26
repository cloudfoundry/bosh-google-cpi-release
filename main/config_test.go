package main_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"

	. "github.com/frodenas/bosh-google-cpi/main"

	bgcaction "github.com/frodenas/bosh-google-cpi/action"

	"github.com/frodenas/bosh-google-cpi/registry"
)

var validGoogleConfig = GoogleConfig{
	Project:         "fake-project",
	JsonKey:         "{}",
	DefaultZone:     "fake-default-zone",
	AccessKeyId:     "fake-access-key-id",
	SecretAccessKey: "fake-secret-access-key",
}

var validActionsOptions = bgcaction.ConcreteFactoryOptions{
	Agent: registry.AgentOptions{
		Mbus: "fake-mbus",
		Ntp:  []string{},
		Blobstore: registry.BlobstoreOptions{
			Type: "fake-blobstore-type",
		},
	},
	Registry: registry.RegistryOptions{
		Schema:   "http",
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
			config.Google.Project = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Google configuration"))
		})

		It("returns error if actions section is not valid", func() {
			config.Actions.Agent = registry.AgentOptions{}
			config.Actions.Registry = registry.RegistryOptions{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Actions configuration"))
		})
	})
})

var _ = Describe("GoogleConfig", func() {
	var (
		config GoogleConfig
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validGoogleConfig
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

		It("returns error if JsonKey is empty", func() {
			config.JsonKey = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty JsonKey"))
		})

		It("returns error if DefaultZone is empty", func() {
			config.DefaultZone = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty DefaultZone"))
		})
		It("returns error if AccessKeyId is empty", func() {
			config.AccessKeyId = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty AccessKeyId"))
		})

		It("returns error if SecretAccessKey is empty", func() {
			config.SecretAccessKey = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty SecretAccessKey"))
		})
	})
})
