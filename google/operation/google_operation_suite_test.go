package goperation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleOperationService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Operation Service Suite")
}
