package image_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestImageService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Image Service Suite")
}
