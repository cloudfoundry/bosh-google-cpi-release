package server_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRegistryServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Server Suite")
}
