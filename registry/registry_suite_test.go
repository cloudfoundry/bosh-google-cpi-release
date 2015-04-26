package registry_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRegistryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Service Suite")
}
