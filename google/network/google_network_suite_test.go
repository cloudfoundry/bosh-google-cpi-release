package gnetwork

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleNetworkService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Network Service Suite")
}
