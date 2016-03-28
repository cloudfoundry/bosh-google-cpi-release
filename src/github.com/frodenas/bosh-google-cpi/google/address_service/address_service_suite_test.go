package address_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAddressService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Address Service Suite")
}
