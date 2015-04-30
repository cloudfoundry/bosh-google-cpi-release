package store_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRegistryStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Store Suite")
}
