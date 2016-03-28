package targetpool_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTargetPoolService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Target Pool Service Suite")
}
