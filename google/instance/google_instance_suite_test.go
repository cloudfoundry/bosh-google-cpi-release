package ginstance_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleInstanceService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Instance Service Suite")
}
