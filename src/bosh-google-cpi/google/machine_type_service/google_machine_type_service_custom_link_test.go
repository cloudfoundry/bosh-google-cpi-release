package machinetype_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/google/machine_type_service"
)

var _ = Describe("GoogleMachineTypeServiceCustomLink", func() {
	var subject GoogleMachineTypeService

	BeforeEach(func() {
		subject = NewGoogleMachineTypeService("foo", nil, nil)
	})

	It("provides a default N1 sized machine type link", func() {
		Expect(subject.CustomLink(2, 2048, "us-east1-d", "")).To(Equal("https://www.googleapis.com/compute/v1/projects/foo/zones/us-east1-d/machineTypes/custom-2-2048"))
	})

	It("provides a explict N1 sized machine type link", func() {
		Expect(subject.CustomLink(2, 2048, "us-east1-d", "N1")).To(Equal("https://www.googleapis.com/compute/v1/projects/foo/zones/us-east1-d/machineTypes/custom-2-2048"))
	})

	It("provides a normal E2 sized machine type link", func() {
		Expect(subject.CustomLink(2, 2048, "us-east1-d", "E2")).To(Equal("https://www.googleapis.com/compute/v1/projects/foo/zones/us-east1-d/machineTypes/e2-custom-2-2048"))
	})

	It("provides an extended machine type link", func() {
		Expect(subject.CustomLink(4, 26880, "us-east1-d", "")).To(Equal("https://www.googleapis.com/compute/v1/projects/foo/zones/us-east1-d/machineTypes/custom-4-26880-ext"))
	})
})
