package gtargetpool

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleTargetPoolService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Target Pool Service Suite")
}
