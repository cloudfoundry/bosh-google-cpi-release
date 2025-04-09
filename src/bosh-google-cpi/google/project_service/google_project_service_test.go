package project_test

import (
	project "bosh-google-cpi/google/project_service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoogleProjectService", func() {
	var (
		service project.GoogleProjectService
	)

	BeforeEach(func() {
		service = project.GoogleProjectService{
			ProjectId: "foo",
		}
	})

	It("uses the default when none is specified", func() {
		Expect(service.Find("")).To(Equal("foo"))
	})

	It("prefers the caller provided project", func() {
		Expect(service.Find("bar")).To(Equal("bar"))
	})
})
