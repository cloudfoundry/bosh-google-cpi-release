package gimage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoogleImageService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Image Service Suite")
}
