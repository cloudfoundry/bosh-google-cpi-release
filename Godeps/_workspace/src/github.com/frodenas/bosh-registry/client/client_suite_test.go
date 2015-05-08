package registry_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRegistryClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Client Suite")
}
