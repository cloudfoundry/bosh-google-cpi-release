package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	"github.com/frodenas/bosh-registry/client"
)

var _ = Describe("ConcreteFactoryOptions", func() {
	var (
		options ConcreteFactoryOptions

		validOptions = ConcreteFactoryOptions{
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
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			options = validOptions
		})

		It("does not return error if all fields are valid", func() {
			err := options.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if agent section is not valid", func() {
			options.Agent = registry.AgentOptions{}

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Agent configuration"))
		})

		It("returns error if registry section is not valid", func() {
			options.Registry = registry.ClientOptions{}

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Registry configuration"))
		})
	})
})
