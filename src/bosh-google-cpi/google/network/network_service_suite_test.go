package network_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNetworkService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Network Service Suite")
}
