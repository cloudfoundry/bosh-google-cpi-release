package gaddress

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleAddressService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Address Service Suite")
}
