package client_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Client Suite")
}
