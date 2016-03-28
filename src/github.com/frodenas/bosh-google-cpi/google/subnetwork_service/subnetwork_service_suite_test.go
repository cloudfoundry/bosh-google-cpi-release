package subnetwork_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSubnetworkService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subnetwork Service Suite")
}
